package rule

import (
	"net/http"
	"testing"

	"github.com/rocinax/rigis/pkg/rule"
)

func TestOrRule(t *testing.T) {

	var trueTrueRules []rule.Rule
	trueTrueRules = append(trueTrueRules, rule.NewTrueRule())
	trueTrueRules = append(trueTrueRules, rule.NewTrueRule())

	var trueFalseRules []rule.Rule
	trueFalseRules = append(trueFalseRules, rule.NewTrueRule())
	trueFalseRules = append(trueFalseRules, rule.NewFalseRule())

	var falseFalseRules []rule.Rule
	falseFalseRules = append(falseFalseRules, rule.NewFalseRule())
	falseFalseRules = append(falseFalseRules, rule.NewFalseRule())

	if rule.NewOrRule(trueTrueRules).Execute(&http.Request{}) == false {
		t.Error("OrRule{Child1: true, Child2: true")
	}

	if rule.NewOrRule(trueFalseRules).Execute(&http.Request{}) == false {
		t.Error("OrRule{Child1: true, Child2: false")
	}

	if rule.NewOrRule(falseFalseRules).Execute(&http.Request{}) == true {
		t.Error("OrRule{Child1: false, Child2: false")
	}

}
