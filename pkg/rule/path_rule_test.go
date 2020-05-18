package rule

import (
	"net/http"
	"testing"
)

type pathRuleTestConfig struct {
	path       string
	exlude     bool
	testPath   string
	testResult bool
}

func TestPathRule(t *testing.T) {

	var pathRuleTestConfigs []pathRuleTestConfig
	pathRuleTestConfigs = append(pathRuleTestConfigs, pathRuleTestConfig{
		path:       "/rocinax",
		exlude:     false,
		testPath:   "/rocinax",
		testResult: true,
	})
	pathRuleTestConfigs = append(pathRuleTestConfigs, pathRuleTestConfig{
		path:       "/rocinax",
		exlude:     true,
		testPath:   "/rocinax",
		testResult: false,
	})
	pathRuleTestConfigs = append(pathRuleTestConfigs, pathRuleTestConfig{
		path:       "/rocinax",
		exlude:     false,
		testPath:   "/rocinax/rigis",
		testResult: true,
	})
	pathRuleTestConfigs = append(pathRuleTestConfigs, pathRuleTestConfig{
		path:       "/rocinax",
		exlude:     true,
		testPath:   "/rocinax/rigis",
		testResult: false,
	})
	pathRuleTestConfigs = append(pathRuleTestConfigs, pathRuleTestConfig{
		path:       "/rocinax",
		exlude:     false,
		testPath:   "/ROCINAX/RIGIS",
		testResult: true,
	})
	pathRuleTestConfigs = append(pathRuleTestConfigs, pathRuleTestConfig{
		path:       "/rocinax",
		exlude:     true,
		testPath:   "/ROCINAX/RIGIS",
		testResult: false,
	})

	for i := 0; i < len(pathRuleTestConfigs); i++ {
		pathRule := NewPathRule(
			PathSetting{
				Path:   pathRuleTestConfigs[i].path,
				Exlude: pathRuleTestConfigs[i].exlude,
			},
		)

		if pathRule.Execute(&http.Request{Host: pathRuleTestConfigs[i].testPath}) != pathRuleTestConfigs[i].testResult {
			t.Errorf(
				"PathRule{Path: %s, Exlude: %t -> RequestPath: %s ",
				pathRuleTestConfigs[i].path,
				pathRuleTestConfigs[i].exlude,
				pathRuleTestConfigs[i].testPath,
			)
		}
	}
}
