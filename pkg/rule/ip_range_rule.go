package rule

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// IPRangeSetting Geolocation Country Rule Setting
type IPRangeSetting struct {
	CIDR    string
	Exclude bool
}

type ipRangeRule struct {
	cidr    *net.IPNet
	exclude bool
}

// NewIPRangeRule :
func NewIPRangeRule(setting IPRangeSetting) Rule {

	_, ipnet, err := net.ParseCIDR(setting.CIDR)
	if err != nil {
		panic(errors.New("IPRangeRule:InvalidCIDRFormat"))
	}

	return ipRangeRule{
		cidr:    ipnet,
		exclude: setting.Exclude,
	}
}

// Execute Execute IP Range Based Rule
func (r ipRangeRule) Execute(req *http.Request) bool {

	remoteIP := net.ParseIP(strings.Split(req.RemoteAddr, ":")[0])

	if r.cidr.Contains(remoteIP) == !r.exclude {
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Tracef("IPRangeRule:Success:%s", remoteIP)
		return true
	}

	logrus.WithFields(logrus.Fields{
		"type": "rule",
		"app":  "rigis",
	}).Tracef("IPRangeRule:Failed:%s", remoteIP)
	return false
}
