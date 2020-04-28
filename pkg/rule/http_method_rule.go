package rule

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// HTTPMethodSetting HTTPMethodRule Setting
type HTTPMethodSetting struct {
	Method  string
	Exclude bool
}

type httpMethodRule struct {
	method  string
	exclude bool
}

// NewHTTPMethodRule :
func NewHTTPMethodRule(setting HTTPMethodSetting) Rule {
	return httpMethodRule{
		method:  setting.Method,
		exclude: setting.Exclude,
	}
}

// Execute Execute HTTP Method Based Rule
func (r httpMethodRule) Execute(req *http.Request) bool {

	if strings.ToUpper(r.method) == strings.ToUpper(req.Method) && !r.exclude {
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Tracef("HTTPMethodRule:Success:%s", r.method)
		return true
	}

	logrus.WithFields(logrus.Fields{
		"type": "rule",
		"app":  "rigis",
	}).Tracef("HTTPMethodRule:Failed:%s", r.method)
	return false
}
