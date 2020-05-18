package rule

import (
	"net/http"
	"testing"
)

func TestTrueRule(t *testing.T) {

	if NewTrueRule().Execute(&http.Request{}) == false {
		t.Error("TrueRule Failed")
	}
}
