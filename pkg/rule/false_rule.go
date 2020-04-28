package rule

import (
	"net/http"
)

type falseRule struct{}

// NewFalseRule :
func NewFalseRule() Rule {
	return falseRule{}
}

// Execute Execute Always False Rule
func (r falseRule) Execute(req *http.Request) bool {
	return false
}
