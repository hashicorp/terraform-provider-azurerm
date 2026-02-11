package databaselink

func HasDbToLink(from, to []string) bool {
	fromMap := make(map[string]bool, len(from))
	for _, id := range from {
		fromMap[id] = true
	}

	for _, id := range to {
		if !fromMap[id] {
			return true
		}
	}
	return false
}
