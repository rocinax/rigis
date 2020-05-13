package rigis

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"text/template"

	"github.com/rocinax/rigis/pkg/rule"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Rigis Rocinax Rigis
type Rigis interface {
	ServeHTTP(res http.ResponseWriter, req *http.Request)
}

type rigis struct {
	filter filter
	nodes  []node
}

// NewRigis :
func NewRigis(config Config) Rigis {

	var nodes []node

	for i := 0; i < len(config.Nodes); i++ {
		nodes = append(nodes, newNode(
			parseRule(config.Nodes[i].Rule),
			newFilter(
				parseRule(config.Nodes[i].Filter.Rule),
				config.Nodes[i].Filter.Accept,
			),
			config.Nodes[i].BalanceType,
			newBackendHosts(config.Nodes[i].BackendHosts),
		))
	}

	return rigis{
		filter: newFilter(
			parseRule(config.Filter.Rule),
			config.Filter.Accept,
		),
		nodes: nodes,
	}
}

func parseRule(configRule ConfigRule) rule.Rule {

	switch configRule.Type {
	case "AndRule":
		var rules []rule.Rule
		for i := 0; i < len(configRule.Rules); i++ {
			rules = append(rules, parseRule(configRule.Rules[i]))
		}
		return rule.NewAndRule(rules)

	case "OrRule":
		var rules []rule.Rule
		for i := 0; i < len(configRule.Rules); i++ {
			rules = append(rules, parseRule(configRule.Rules[i]))
		}
		return rule.NewOrRule(rules)

	default:
		return rule.NewRule(configRule.Type, configRule.Setting)
	}
}

func newBackendHosts(configHosts []ConfigBackendHost) []backendHost {

	var backendHosts []backendHost

	for i := 0; i < len(configHosts); i++ {

		beURL, err := url.Parse(configHosts[i].URL)

		if err != nil {
			// FIXME panic message
			panic(errors.New("BackendHost.URL: invalid url format"))
		}

		backendHosts = append(backendHosts, newBackendHost(
			configHosts[i].Weight,
			beURL,
		))
	}

	return backendHosts
}

// ServeHTTP Custrum HTTP Handle Function
func (r rigis) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// if filter match then bad request
	if !r.executeFilter(req) {
		logrus.WithFields(FormatErrorLog(req)).Error("Request Forbidden.")
		ResponseError(rw, req, map[string]interface{}{
			"Status":           http.StatusForbidden,
			"Error":            "Forbidden",
			"ErrorDescription": "Request Forbidden.",
		})
		return
	}

	for i := 0; i < len(r.nodes); i++ {
		// if rule unmatch then skip this node
		if !r.nodes[i].executeRule(req) {
			continue
		}

		// if filter match then bad request
		if !r.nodes[i].executeFilter(req) {
			logrus.WithFields(FormatErrorLog(req)).Error("Request Forbidden.")
			ResponseError(rw, req, map[string]interface{}{
				"Status":           http.StatusForbidden,
				"Error":            "Forbidden",
				"ErrorDescription": "Request Forbidden.",
			})
			return
		}

		backend := r.nodes[i].getBackend()
		backend.ServeHTTP(rw, req)

		return
	}

	// not match all rule
	logrus.WithFields(FormatErrorLog(req)).Error("Contnts Not Found.")
	ResponseError(rw, req, map[string]interface{}{
		"Status":           http.StatusNotFound,
		"Error":            "Not Found",
		"ErrorDescription": "Contents Not Found.",
	})
}

func (r rigis) executeFilter(req *http.Request) bool {
	return r.filter.Execute(req)
}

// ResponseError :
func ResponseError(rw http.ResponseWriter, req *http.Request, data map[string]interface{}) {

	// http cache control header
	rw.Header().Set("Pragma", "no-cache")
	rw.Header().Set("Expires", "-1")
	rw.Header().Set("Cache-Control", "no-cache,no-store")

	// http security control header
	rw.Header().Set("Server", viper.GetString("ServerName"))
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("X-Frame-Options", "deny")
	rw.Header().Set("X-XSS-Protection", "1; mode=block")
	rw.Header().Set("Referrer-Policy", "no-referrer")
	rw.Header().Set("Content-Security-Policy", "")

	// response data
	rw.WriteHeader(data["Status"].(int))
	errorTemplate := path.Join(viper.GetString("TemplateDir"), viper.GetString("ErrorTemplate"))
	templateHTML := template.Must(template.ParseFiles(errorTemplate))
	templateHTML.Execute(rw, data)
	return
}

// FormatAccessLog :
func FormatAccessLog(res *http.Response) map[string]interface{} {

	var remoteAddr string
	var remotePort string

	if res.Request.Header.Get("X-Real-IP") != "" {
		remoteAddr = res.Request.Header.Get("X-Real-IP")
		remotePort = ""
	} else if res.Request.Header.Get("X-Forwarded-For") != "" {
		remoteAddr = strings.Split(res.Request.Header.Get("X-Forwarded-For"), ",")[0]
		remotePort = ""
	} else {
		remoteAddr = res.Request.URL.Hostname()
		remotePort = res.Request.URL.Port()
	}

	return logrus.Fields{
		"type":            "access",
		"app":             "rigis",
		"remote_address":  remoteAddr,
		"remote_port":     remotePort,
		"host":            res.Request.Host,
		"method":          res.Request.Method,
		"request_uri":     res.Request.RequestURI,
		"protocol":        res.Request.Proto,
		"referer":         res.Request.Referer(),
		"user_agent":      res.Request.UserAgent(),
		"status_code":     res.StatusCode,
		"contents_length": res.ContentLength,
		"user":            fmt.Sprintf("%x", (md5.Sum(([]byte(res.Request.Header.Get("Cookie")))[:]))),
	}
}

// FormatErrorLog :
func FormatErrorLog(req *http.Request) map[string]interface{} {

	var remoteAddr string
	var remotePort string

	if req.Header.Get("X-Read-IP") != "" {
		remoteAddr = req.Header.Get("X-Read-IP")
		remotePort = ""
	} else if req.Header.Get("X-Forwarded-For") != "" {
		remoteAddr = strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0]
		remotePort = ""
	} else {
		remoteAddr = strings.Split(req.RemoteAddr, ":")[0]
		remotePort = strings.Split(req.RemoteAddr, ":")[1]
	}

	return logrus.Fields{
		"type":            "access",
		"app":             "rigis",
		"remote_address":  remoteAddr,
		"remote_port":     remotePort,
		"host":            req.Host,
		"method":          req.Method,
		"request_uri":     req.RequestURI,
		"protocol":        req.Proto,
		"referer":         req.Referer(),
		"user_agent":      req.UserAgent(),
		"status_code":     http.StatusServiceUnavailable,
		"contents_length": 0,
		"user":            fmt.Sprintf("%x", (md5.Sum(([]byte(req.Header.Get("Cookie")))[:]))),
	}
}
