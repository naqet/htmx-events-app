package utils

import "htmx-events-app/db"

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
