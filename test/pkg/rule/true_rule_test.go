package rule

import (
	"net/http"
	"testing"

	"github.com/rocinax/rigis/pkg/rule"
)

func TestTrueRule(t *testing.T) {

	if rule.NewTrueRule().Execute(&http.Request{}) == false {
		t.Error("TrueRule Failed")
	}
}
