package rule

import (
	"net/http"
	"testing"
)

func TestOrRule(t *testing.T) {

	var trueTrueRules []Rule
	trueTrueRules = append(trueTrueRules, NewTrueRule())
	trueTrueRules = append(trueTrueRules, NewTrueRule())

	var trueFalseRules []Rule
	trueFalseRules = append(trueFalseRules, NewTrueRule())
	trueFalseRules = append(trueFalseRules, NewFalseRule())

	var falseFalseRules []Rule
	falseFalseRules = append(falseFalseRules, NewFalseRule())
	falseFalseRules = append(falseFalseRules, NewFalseRule())

	if NewOrRule(trueTrueRules).Execute(&http.Request{}) == false {
		t.Error("OrRule{Child1: true, Child2: true")
	}

	if NewOrRule(trueFalseRules).Execute(&http.Request{}) == false {
		t.Error("OrRule{Child1: true, Child2: false")
	}

	if NewOrRule(falseFalseRules).Execute(&http.Request{}) == true {
		t.Error("OrRule{Child1: false, Child2: false")
	}

}
