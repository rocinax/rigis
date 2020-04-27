package rule

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type orRule struct {
	rules []Rule
}

// NewOrRule :
func NewOrRule(rules []Rule) Rule {
	return orRule{
		rules: rules,
	}
}

// Execute Execute And Rule
func (r orRule) Execute(req *http.Request) bool {

	logrus.WithFields(logrus.Fields{
		"type": "rule",
		"app":  "rigis",
	}).Trace("OrRule:Start")

	if len(r.rules) == 0 {
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Trace("OrRule:Failed")
		return false
	}

	for i := 0; i < len(r.rules); i++ {
		if r.rules[i].Execute(req) {
			logrus.WithFields(logrus.Fields{
				"type": "rule",
				"app":  "rigis",
			}).Trace("OrRule:Success")
			return true
		}
	}

	logrus.WithFields(logrus.Fields{
		"type": "rule",
		"app":  "rigis",
	}).Trace("OrRule:Failed")
	return false
}
