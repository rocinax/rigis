package rule

import (
	"net/http"
	"testing"

	"github.com/rocinax/rigis/pkg/rule"
)

func TestTimeRule(t *testing.T) {

	includeTimeRule := rule.NewTimeRule(
		rule.TimeSetting{
			StartHour:    0,
			StartMinutes: 0,
			EndHour:      23,
			EndMinutes:   59,
			Exclude:      false,
		},
	)

	excludeTimeRule := rule.NewTimeRule(
		rule.TimeSetting{
			StartHour:    0,
			StartMinutes: 0,
			EndHour:      23,
			EndMinutes:   59,
			Exclude:      true,
		},
	)

	if includeTimeRule.Execute(&http.Request{}) == false {
		t.Error("TimeRule{Exlude: false")
	}

	if excludeTimeRule.Execute(&http.Request{}) == true {
		t.Error("TimeRule{Exlude: ture")
	}
}
