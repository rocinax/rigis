package rigis

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/rocinax/rigis/pkg/rule"
	"github.com/sirupsen/logrus"
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
		logrus.WithFields(formatErrorLog(req)).Error("Request Forbidden.")
		responseError(rw, req, getErrorEntity(http.StatusForbidden))
		return
	}

	for i := 0; i < len(r.nodes); i++ {
		// if rule unmatch then skip this node
		if !r.nodes[i].executeRule(req) {
			continue
		}

		// if filter match then bad request
		if !r.nodes[i].executeFilter(req) {
			logrus.WithFields(formatErrorLog(req)).Error("Request Forbidden.")
			responseError(rw, req, getErrorEntity(http.StatusForbidden))
			return
		}

		r.nodes[i].serveHTTP(rw, req)
		return
	}

	// not match all rule
	logrus.WithFields(formatErrorLog(req)).Error("Contnts Not Found.")
	responseError(rw, req, getErrorEntity(http.StatusNotFound))
}

func (r rigis) executeFilter(req *http.Request) bool {
	return r.filter.Execute(req)
}
