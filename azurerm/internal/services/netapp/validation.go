package netapp

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func ValidateNetAppAccountName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[-_\da-zA-Z]{3,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppPoolName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z][-_\da-zA-Z]{2,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and start with letters or numbers and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppVolumeName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z][-_\da-zA-Z]{0,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length and start with letters and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppVolumeVolumePath(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z][-\da-zA-Z]{0,79}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 80 characters in length and start with letters and contains only letters, numbers or hyphens.", k))
	}

	return warnings, errors
}

func ValidateNetAppSnapshotName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z][-_\da-zA-Z]{3,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 4 and 64 characters in length and start with letters or numbers and contains only letters, numbers, underscore or hyphens.", k))
	}

	return warnings, errors
}

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
