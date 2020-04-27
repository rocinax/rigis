package rule

import (
	"net/http"
	"testing"

	"github.com/rocinax/rigis/pkg/rule"
)

type httpMethodRuleTestConfig struct {
	method     string
	exclude    bool
	testMethod string
	testResult bool
}

func TestHTTPpMethodRule(t *testing.T) {

	var httpMethodRuleTestConfigs []httpMethodRuleTestConfig
	httpMethodRuleTestConfigs = append(httpMethodRuleTestConfigs, httpMethodRuleTestConfig{
		method:     "get",
		exclude:    false,
		testMethod: "get",
		testResult: true,
	})
	httpMethodRuleTestConfigs = append(httpMethodRuleTestConfigs, httpMethodRuleTestConfig{
		method:     "get",
		exclude:    true,
		testMethod: "get",
		testResult: false,
	})
	httpMethodRuleTestConfigs = append(httpMethodRuleTestConfigs, httpMethodRuleTestConfig{
		method:     "GET",
		exclude:    false,
		testMethod: "get",
		testResult: true,
	})
	httpMethodRuleTestConfigs = append(httpMethodRuleTestConfigs, httpMethodRuleTestConfig{
		method:     "GET",
		exclude:    true,
		testMethod: "get",
		testResult: false,
	})

	for i := 0; i < len(httpMethodRuleTestConfigs); i++ {
		httpMethodRule := rule.NewHTTPMethodRule(
			rule.HTTPMethodSetting{
				Method:  httpMethodRuleTestConfigs[i].method,
				Exclude: httpMethodRuleTestConfigs[i].exclude,
			},
		)

		if httpMethodRule.Execute(&http.Request{Method: httpMethodRuleTestConfigs[i].testMethod}) != httpMethodRuleTestConfigs[i].testResult {
			t.Errorf(
				"HTTPMethodRule{Method: %s, Exclude: %t -> RequestMethod: %s ",
				httpMethodRuleTestConfigs[i].method,
				httpMethodRuleTestConfigs[i].exclude,
				httpMethodRuleTestConfigs[i].testMethod,
			)
		}
	}
}
