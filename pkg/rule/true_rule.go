package rule

import (
	"net/http"
)

type trueRule struct{}

// NewTrueRule :
func NewTrueRule() Rule {
	return trueRule{}
}

// Execute Execute Always True Rule
func (r trueRule) Execute(req *http.Request) bool {
	return true
}
