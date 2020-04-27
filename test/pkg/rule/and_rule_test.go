package rule

import (
	"net/http"
	"testing"

	"github.com/rocinax/rigis/pkg/rule"
)

func TestAndRule(t *testing.T) {

	var trueTrueRules []rule.Rule
	trueTrueRules = append(trueTrueRules, rule.NewTrueRule())
	trueTrueRules = append(trueTrueRules, rule.NewTrueRule())

	var trueFalseRules []rule.Rule
	trueFalseRules = append(trueFalseRules, rule.NewTrueRule())
	trueFalseRules = append(trueFalseRules, rule.NewFalseRule())

	var falseFalseRules []rule.Rule
	falseFalseRules = append(falseFalseRules, rule.NewFalseRule())
	falseFalseRules = append(falseFalseRules, rule.NewFalseRule())

	if rule.NewAndRule(trueTrueRules).Execute(&http.Request{}) == false {
		t.Error("AndRule{Child1: true, Child2: true")
	}

	if rule.NewAndRule(trueFalseRules).Execute(&http.Request{}) == true {
		t.Error("AndRule{Child1: true, Child2: false")
	}

	if rule.NewAndRule(falseFalseRules).Execute(&http.Request{}) == true {
		t.Error("AndRule{Child1: false, Child2: false")
	}

}
