package vevents;

func getFirstLetter(name string) string {
	if len(name) == 0 {
		return ""
	}
	return string(name[0])
}

