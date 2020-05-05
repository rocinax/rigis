package rule

import (
	"encoding/json"
	"net/http"
)

// Rule Balancind Rule
type Rule interface {
	Execute(req *http.Request) bool
}

// NewRule :
func NewRule(ruleType string, setting map[string]interface{}) Rule {
	switch ruleType {
	case "DomainRule":
		var domainSetting DomainSetting
		mapToStruct(setting, &domainSetting)
		return NewDomainRule(domainSetting)
	case "URIPathRule":
		var pathSetting PathSetting
		mapToStruct(setting, &pathSetting)
		return NewPathRule(pathSetting)
	case "HTTPMethodRule":
		var httpMethodSetting HTTPMethodSetting
		mapToStruct(setting, &httpMethodSetting)
		return NewHTTPMethodRule(httpMethodSetting)
	case "IPRangeRule":
		var ipRangeSetting IPRangeSetting
		mapToStruct(setting, &ipRangeSetting)
		return NewIPRangeRule(ipRangeSetting)
	case "TimeRule":
		var timeSetting TimeSetting
		mapToStruct(setting, &timeSetting)
		return NewTimeRule(timeSetting)
	case "TrueRule":
		return NewTrueRule()
	case "FalseRule":
		return NewFalseRule()
	default:
		return NewFalseRule()
	}
}

func mapToStruct(m map[string]interface{}, val interface{}) {
	tmp, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		panic(err)
	}
}
