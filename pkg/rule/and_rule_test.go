package rule

import (
	"net/http"
	"testing"
)

func TestAndRule(t *testing.T) {

	var trueTrueRules []Rule
	trueTrueRules = append(trueTrueRules, NewTrueRule())
	trueTrueRules = append(trueTrueRules, NewTrueRule())

	var trueFalseRules []Rule
	trueFalseRules = append(trueFalseRules, NewTrueRule())
	trueFalseRules = append(trueFalseRules, NewFalseRule())

	var falseFalseRules []Rule
	falseFalseRules = append(falseFalseRules, NewFalseRule())
	falseFalseRules = append(falseFalseRules, NewFalseRule())

	if NewAndRule(trueTrueRules).Execute(&http.Request{}) == false {
		t.Error("AndRule{Child1: true, Child2: true")
	}

	if NewAndRule(trueFalseRules).Execute(&http.Request{}) == true {
		t.Error("AndRule{Child1: true, Child2: false")
	}

	if NewAndRule(falseFalseRules).Execute(&http.Request{}) == true {
		t.Error("AndRule{Child1: false, Child2: false")
	}
}
