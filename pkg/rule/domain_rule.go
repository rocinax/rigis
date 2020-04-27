package rule

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// DomainSetting Domain Rule Setting
type DomainSetting struct {
	Domain           string
	IncludeSubDomain bool
}

type domainRule struct {
	domain           string
	includeSubDomain bool
}

// NewDomainRule :
func NewDomainRule(setting DomainSetting) Rule {
	return domainRule{
		domain:           setting.Domain,
		includeSubDomain: setting.IncludeSubDomain,
	}
}

// Execute Execute Domain Based Rule
func (r domainRule) Execute(req *http.Request) bool {

	if strings.ToUpper(r.domain) == strings.ToUpper(req.Host) {
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Tracef("DomainRule:Success:%s", r.domain)
		return true
	}

	if r.includeSubDomain && strings.HasSuffix(strings.ToUpper(req.Host), "."+strings.ToUpper(r.domain)) {
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Tracef("DomainRule:Success:%s", r.domain)
		return true
	}

	logrus.WithFields(logrus.Fields{
		"type": "rule",
		"app":  "rigis",
	}).Tracef("DomainRule:Failed:%s", r.domain)
	return false
}
