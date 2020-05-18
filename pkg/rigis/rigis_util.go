package rigis

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"path"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func normalizeRpPath(stringA, stringB string) string {
	aslash := strings.HasSuffix(stringA, "/")
	bslash := strings.HasPrefix(stringB, "/")
	switch {
	case aslash && bslash:
		return stringA + stringB[1:]
	case !aslash && !bslash:
		return stringA + "/" + stringB
	}
	return stringA + stringB
}

func normalizeRawQuery(rpRawQuery string, reqRawQuery string) string {

	if rpRawQuery == "" || reqRawQuery == "" {
		return rpRawQuery + reqRawQuery
	}

	return rpRawQuery + "&" + reqRawQuery
}

func getXForwardedProto(headerProto string, tlsEnable bool) string {

	if headerProto == "http" || headerProto == "https" {
		return headerProto
	}

	if tlsEnable {
		return "https"
	}
	return "http"
}

func getXForwardedHost(xForwardedHost string, host string) string {

	if xForwardedHost != "" {
		return xForwardedHost
	}

	return host
}

func getUserAgent(userAgent string) string {
	if userAgent == "" {
		return userAgent
	}

	return viper.GetString("ProxyUserAgent")
}

func getRemoteAddr(req *http.Request) string {

	if req.Header.Get("X-Real-IP") != "" {
		return req.Header.Get("X-Real-IP")
	}

	if req.Header.Get("X-Forwarded-For") != "" {
		return strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0]
	}

	return req.RemoteAddr[strings.LastIndex(req.RemoteAddr, ":"):len(req.RemoteAddr)]
}

func getErrorEntity(httpStatus int) map[string]interface{} {

	switch httpStatus {
	case http.StatusForbidden:
		return map[string]interface{}{
			"Status":           http.StatusForbidden,
			"Error":            "Forbidden",
			"ErrorDescription": "Request Forbidden.",
		}
	case http.StatusNotFound:
		return map[string]interface{}{
			"Status":           http.StatusNotFound,
			"Error":            "Not Found",
			"ErrorDescription": "Contents Not Found.",
		}
	case http.StatusServiceUnavailable:
		return map[string]interface{}{
			"Status":           http.StatusServiceUnavailable,
			"Error":            "Service Unavailable",
			"ErrorDescription": "HTTP Service Unavailable.",
		}
	}

	return map[string]interface{}{
		"Status":           http.StatusServiceUnavailable,
		"Error":            "Service Unavailable",
		"ErrorDescription": "HTTP Service Unavailable.",
	}
}

func formatAccessLog(res *http.Response) map[string]interface{} {

	return logrus.Fields{
		"type":            "access",
		"app":             "rigis",
		"remote_address":  getRemoteAddr(res.Request),
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

func formatErrorLog(req *http.Request) map[string]interface{} {

	return logrus.Fields{
		"type":            "access",
		"app":             "rigis",
		"remote_address":  getRemoteAddr(req),
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

func responseError(rw http.ResponseWriter, req *http.Request, data map[string]interface{}) {

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
