package common

// Remove duplicates from the input array and return unify array (without duplicated elements)
func RemoveDuplicatesFromStringArray(elements []string) []string {
	visited := map[string]bool{}
	result := []string{}

	for v := range elements {
		if !visited[elements[v]] {
			visited[elements[v]] = true          // Mark the element as visited.
			result = append(result, elements[v]) // Add it to the result.
		}
	}

	return result
}
