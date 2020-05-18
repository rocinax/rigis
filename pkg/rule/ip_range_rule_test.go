package rule

import (
	"net/http"
	"testing"
)

type ipRangeRuleTestConfig struct {
	cidr           string
	exclude        bool
	testRemoteAddr string
	testResult     bool
}

func TestIPRangeRule(t *testing.T) {

	var ipRangeRuleTestConfigs []ipRangeRuleTestConfig
	ipRangeRuleTestConfigs = append(ipRangeRuleTestConfigs, ipRangeRuleTestConfig{
		cidr:           "192.168.0.0/16",
		exclude:        false,
		testRemoteAddr: "192.168.1.1:50000",
		testResult:     true,
	})
	ipRangeRuleTestConfigs = append(ipRangeRuleTestConfigs, ipRangeRuleTestConfig{
		cidr:           "192.168.0.0/16",
		exclude:        true,
		testRemoteAddr: "192.168.1.1:50000",
		testResult:     false,
	})
	ipRangeRuleTestConfigs = append(ipRangeRuleTestConfigs, ipRangeRuleTestConfig{
		cidr:           "192.168.1.1/32",
		exclude:        false,
		testRemoteAddr: "192.168.1.1:50000",
		testResult:     true,
	})
	ipRangeRuleTestConfigs = append(ipRangeRuleTestConfigs, ipRangeRuleTestConfig{
		cidr:           "192.168.1.1/32",
		exclude:        true,
		testRemoteAddr: "192.168.1.1:50000",
		testResult:     false,
	})

	for i := 0; i < len(ipRangeRuleTestConfigs); i++ {
		ipRangeRule := NewIPRangeRule(
			IPRangeSetting{
				CIDR:    ipRangeRuleTestConfigs[i].cidr,
				Exclude: ipRangeRuleTestConfigs[i].exclude,
			},
		)

		if ipRangeRule.Execute(&http.Request{RemoteAddr: ipRangeRuleTestConfigs[i].testRemoteAddr}) != ipRangeRuleTestConfigs[i].testResult {
			t.Errorf(
				"IPRangeRule{CIDR: %s, Exclude: %t -> RequestRemoteAddr: %s ",
				ipRangeRuleTestConfigs[i].cidr,
				ipRangeRuleTestConfigs[i].exclude,
				ipRangeRuleTestConfigs[i].testRemoteAddr,
			)
		}
	}
}
