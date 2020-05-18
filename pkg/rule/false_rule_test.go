package rule

import (
	"net/http"
	"testing"
)

func TestFalseRule(t *testing.T) {

	if NewFalseRule().Execute(&http.Request{}) == true {
		t.Error("FalseRule Failed")
	}
}
