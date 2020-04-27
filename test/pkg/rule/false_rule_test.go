package rule

import (
	"net/http"
	"testing"

	"github.com/rocinax/rigis/pkg/rule"
)

func TestFalseRule(t *testing.T) {

	if rule.NewFalseRule().Execute(&http.Request{}) == true {
		t.Error("FalseRule Failed")
	}
}
