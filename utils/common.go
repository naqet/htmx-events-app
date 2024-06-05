package utils

import (
	"fmt"
	"htmx-events-app/db"
	"time"
)

func IsEventOwner(email string, event db.Event) bool {
	var isOwner bool

	for _, host := range event.Hosts {
		if host.Email == email {
			isOwner = true
			break
		}
	}

	return isOwner
}

func InTimeSpan(start, end, check time.Time) bool {
	_end := end
	_check := check
	if end.Before(start) {
		_end = end.Add(24 * time.Hour)
		if check.Before(start) {
			_check = check.Add(24 * time.Hour)
		}
	}
	return (_check.Equal(start) || _check.After(start)) && (_check.Equal(_end) || _check.Before(_end))
}

type AgendaSections = map[string][]db.AgendaPoint

func OrganizeAgendaPoints(points []db.AgendaPoint) AgendaSections {
	sections := AgendaSections{}

	for _, point := range points {
		date := fmt.Sprintf("%d %s", point.StartTime.Day(), point.StartTime.Month())
		if _, ok := sections[date]; !ok {
			sections[date] = []db.AgendaPoint{point}
		} else {
			sections[date] = append(sections[date], point)
		}
	}

	return sections
}
