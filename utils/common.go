package utils

import (
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
	return _check.After(start) && _check.Before(_end)
}
