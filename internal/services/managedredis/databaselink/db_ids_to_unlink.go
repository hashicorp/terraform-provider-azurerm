package databaselink

func DbIdsToUnlink(from, to []string) []string {
	result := make([]string, 0)

	toMap := make(map[string]bool, len(to))
	for _, id := range to {
		toMap[id] = true
	}

	for _, id := range from {
		if !toMap[id] {
			result = append(result, id)
		}
	}
	return result
}
