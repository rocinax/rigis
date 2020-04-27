package rule

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// PathSetting URI Path Rule Setting
type PathSetting struct {
	Path   string
	Exlude bool
}

type pathRule struct {
	path   string
	exlude bool
}

// NewPathRule :
func NewPathRule(setting PathSetting) Rule {
	return pathRule{
		path:   setting.Path,
		exlude: setting.Exlude,
	}
}

// Execute Execute URI Path Based Rule
func (r pathRule) Execute(req *http.Request) bool {

	if strings.HasPrefix(strings.ToUpper(r.path), strings.ToUpper(req.RequestURI)) && !r.exlude {
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Tracef("PathRule:Success:%s", req.RequestURI)
		return true
	}

	logrus.WithFields(logrus.Fields{
		"type": "rule",
		"app":  "rigis",
	}).Tracef("PathRule:Failed:%s", req.RequestURI)
	return false
}
