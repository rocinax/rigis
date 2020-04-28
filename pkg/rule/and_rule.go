package rule

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type andRule struct {
	rules []Rule
}

// NewAndRule :
func NewAndRule(rules []Rule) Rule {
	return andRule{
		rules: rules,
	}
}

// Execute Execute And Rule
func (r andRule) Execute(req *http.Request) bool {
	logrus.WithFields(logrus.Fields{
		"type": "rule",
		"app":  "rigis",
	}).Trace("AndRule:Start")

	if len(r.rules) == 0 {
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Trace("AndRule:Failed")
		return false
	}

	for i := 0; i < len(r.rules); i++ {
		if !(r.rules[i].Execute(req)) {
			logrus.WithFields(logrus.Fields{
				"type": "rule",
				"app":  "rigis",
			}).Trace("AndRule:Failed")
			return false
		}
	}

	logrus.WithFields(logrus.Fields{
		"type": "rule",
		"app":  "rigis",
	}).Trace("AndRule:Success")
	return true
}
