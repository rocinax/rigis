package rule

import (
	"net/http"
	"testing"

	"github.com/rocinax/rigis/pkg/rule"
)

type domainRuleTestConfig struct {
	domain           string
	includeSubDomain bool
	testDomain       string
	testResult       bool
}

func TestDomainRule(t *testing.T) {

	var domainRuleTestConfigs []domainRuleTestConfig
	domainRuleTestConfigs = append(domainRuleTestConfigs, domainRuleTestConfig{
		domain:           "rocinax.com",
		includeSubDomain: true,
		testDomain:       "rocinax.com",
		testResult:       true,
	})
	domainRuleTestConfigs = append(domainRuleTestConfigs, domainRuleTestConfig{
		domain:           "rocinax.com",
		includeSubDomain: false,
		testDomain:       "rocinax.com",
		testResult:       true,
	})
	domainRuleTestConfigs = append(domainRuleTestConfigs, domainRuleTestConfig{
		domain:           "rocinax.com",
		includeSubDomain: true,
		testDomain:       "rigis.rocinax.com",
		testResult:       true,
	})
	domainRuleTestConfigs = append(domainRuleTestConfigs, domainRuleTestConfig{
		domain:           "rocinax.com",
		includeSubDomain: false,
		testDomain:       "rigis.rocinax.com",
		testResult:       false,
	})
	domainRuleTestConfigs = append(domainRuleTestConfigs, domainRuleTestConfig{
		domain:           "rocinax.com",
		includeSubDomain: true,
		testDomain:       "rigis.rocinax.com:6443",
		testResult:       false,
	})
	domainRuleTestConfigs = append(domainRuleTestConfigs, domainRuleTestConfig{
		domain:           "rocinax.com",
		includeSubDomain: false,
		testDomain:       "rigis.rocinax.com:6443",
		testResult:       false,
	})
	domainRuleTestConfigs = append(domainRuleTestConfigs, domainRuleTestConfig{
		domain:           "rocinax.com",
		includeSubDomain: true,
		testDomain:       "ROCINAX.COM",
		testResult:       true,
	})
	domainRuleTestConfigs = append(domainRuleTestConfigs, domainRuleTestConfig{
		domain:           "rocinax.com",
		includeSubDomain: false,
		testDomain:       "ROCINAX.COM",
		testResult:       true,
	})

	for i := 0; i < len(domainRuleTestConfigs); i++ {
		domainRule := rule.NewDomainRule(
			rule.DomainSetting{
				Domain:           domainRuleTestConfigs[i].domain,
				IncludeSubDomain: domainRuleTestConfigs[i].includeSubDomain,
			},
		)

		if domainRule.Execute(&http.Request{Host: domainRuleTestConfigs[i].testDomain}) != domainRuleTestConfigs[i].testResult {
			t.Errorf(
				"DomainRule{Domain: %s, IncludeSubdomain: %t -> RequestHost: %s ",
				domainRuleTestConfigs[i].domain,
				domainRuleTestConfigs[i].includeSubDomain,
				domainRuleTestConfigs[i].testDomain,
			)
		}
	}
}
