package rigis

import (
	"net/http"

	"github.com/rocinax/rigis/pkg/rule"
)

type filter struct {
	rule   rule.Rule
	accept bool
}

func newFilter(fRule rule.Rule, accept bool) filter {
	return filter{
		rule:   fRule,
		accept: accept,
	}
}

func (f filter) Execute(req *http.Request) bool {
	return f.rule.Execute(req) == f.accept
}
