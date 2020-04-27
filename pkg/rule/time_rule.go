package rule

import (
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// TimeSetting Time Rule Setting
type TimeSetting struct {
	StartHour    int
	StartMinutes int
	EndHour      int
	EndMinutes   int
	Exclude      bool
}

type timeRule struct {
	ruleDue      int64
	startHour    int
	startMinutes int
	startTime    int64
	endHour      int
	endMinutes   int
	endTime      int64
	exclude      bool
}

// NewTimeRule :
func NewTimeRule(setting TimeSetting) Rule {

	if setting.StartHour < 0 || setting.StartHour > 23 {
		panic(errors.New("TimeRule:OutOfRangeStartTimeHour"))
	}

	if setting.StartMinutes < 0 && setting.StartMinutes > 59 {
		panic(errors.New("TimeRule:OutOfRangeStartTimeMinutes"))
	}

	if setting.EndHour < 0 || setting.EndHour > 23 {
		panic(errors.New("TimeRule:OutOfRangeEndTimeHour"))
	}

	if setting.EndMinutes < 0 && setting.EndMinutes > 59 {
		panic(errors.New("TimeRule:OutOfRangeEndTimeMinutes"))
	}

	if setting.StartHour > setting.EndHour {
		panic(errors.New("TimeRule:EndTimeHourLessThanStartTimeHour"))
	}

	if (setting.StartHour == setting.EndHour) &&
		(setting.StartMinutes > setting.EndMinutes) {
		panic(errors.New("TimeRule:EndTimeMinutesLessThanStartTimeMinutes"))
	}

	return timeRule{
		ruleDue:      time.Now().Unix() - 1,
		startHour:    setting.StartHour,
		startMinutes: setting.StartMinutes,
		endHour:      setting.EndHour,
		endMinutes:   setting.EndMinutes,
		exclude:      setting.Exclude,
	}
}

// Execute Execute Geo Country Based Rule
func (r timeRule) Execute(req *http.Request) bool {

	now := time.Now().Unix()

	// if rule due is over then update time range and rule due
	if now >= r.ruleDue {
		r.ruleDue, r.startTime, r.endTime = r.updateTimeRule()
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Trace("TimeRule:UpdateDue")
	}

	if now >= r.startTime && now <= r.endTime && !r.exclude {
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Trace("TimeRule:Success")
		return true
	}

	if (now <= r.startTime || now >= r.endTime) && r.exclude {
		logrus.WithFields(logrus.Fields{
			"type": "rule",
			"app":  "rigis",
		}).Trace("TimeRule:Success")
		return true
	}

	logrus.WithFields(logrus.Fields{
		"type": "rule",
		"app":  "rigis",
	}).Trace("TimeRule:Failed")
	return false
}

func (r timeRule) updateTimeRule() (int64, int64, int64) {
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)

	// Update Rule Due
	ruleDue := time.Date(
		tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local,
	).Unix()

	start := time.Date(
		today.Year(), today.Month(), today.Day(), r.startHour, r.startMinutes, 0, 0, time.Local,
	).Unix()

	end := time.Date(
		today.Year(), today.Month(), today.Day(), r.endHour, r.endMinutes, 59, 999999999, time.Local,
	).Unix()

	return ruleDue, start, end
}
