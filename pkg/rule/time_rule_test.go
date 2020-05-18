package rule

import (
	"net/http"
	"testing"
)

func TestTimeRule(t *testing.T) {

	includeTimeRule := NewTimeRule(
		TimeSetting{
			StartHour:    0,
			StartMinutes: 0,
			EndHour:      23,
			EndMinutes:   59,
			Exclude:      false,
		},
	)

	excludeTimeRule := NewTimeRule(
		TimeSetting{
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
