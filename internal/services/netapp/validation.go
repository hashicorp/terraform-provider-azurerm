package netapp

import (
	"reflect"
	"strings"
)

// TODO: this should likely be moved to utils, or removed in case of using the values from the source snapshot?

func ValidateSlicesEquality(source, new []string, caseSensitive bool) bool {
	// Fast path
	if len(source) != len(new) {
		return false
	}

	if reflect.DeepEqual(source, new) {
		return true
	}

	// Slow path
	// Source -> New direction
	sourceNewValidatedCount := 0
	for _, sourceItem := range source {
		for _, newItem := range new {
			if caseSensitive {
				if sourceItem == newItem {
					sourceNewValidatedCount++
				}
			} else {
				if strings.EqualFold(sourceItem, newItem) {
					sourceNewValidatedCount++
				}
			}
		}
	}

	// New -> Source direction
	newSourceValidatedCount := 0
	for _, newItem := range source {
		for _, sourceItem := range new {
			if caseSensitive {
				if newItem == sourceItem {
					newSourceValidatedCount++
				}
			} else {
				if strings.EqualFold(newItem, sourceItem) {
					newSourceValidatedCount++
				}
			}
		}
	}

	lengthValidation := sourceNewValidatedCount == len(source) && newSourceValidatedCount == len(source) && sourceNewValidatedCount == len(new) && newSourceValidatedCount == len(new)
	countValidation := sourceNewValidatedCount == newSourceValidatedCount

	return lengthValidation && countValidation
}
