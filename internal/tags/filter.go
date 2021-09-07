package tags

import "strings"

func Filter(tagsMap map[string]*string, tagNames ...string) map[string]*string {
	if len(tagNames) == 0 {
		return tagsMap
	}

	// Build the filter dictionary from passed tag names.
	filterDict := make(map[string]bool)
	for _, name := range tagNames {
		if len(name) > 0 {
			filterDict[strings.ToLower(name)] = true
		}
	}

	// Filter out tag if it exists(case insensitive) in the dictionary.
	tagsRet := make(map[string]*string)
	for k, v := range tagsMap {
		if !filterDict[strings.ToLower(k)] {
			tagsRet[k] = v
		}
	}

	return tagsRet
}
