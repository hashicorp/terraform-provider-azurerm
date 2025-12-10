// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudexadatainfrastructures"
)

func ComputeCount(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return warnings, errors
	}

	if v < 2 || v > 32 {
		errors = append(errors, fmt.Errorf("the compute count must be between %d and %d", 2, 32))
		return warnings, errors
	}

	return warnings, errors
}

func StorageCount(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return warnings, errors
	}

	if v < 3 || v > 64 {
		errors = append(errors, fmt.Errorf("the storage count must be between %d and %d", 3, 64))
		return warnings, errors
	}

	return warnings, errors
}

// MaintenanceWindow validation

func DaysOfWeek(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	validDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	if !slices.Contains(validDaysOfWeek, v) {
		errors = append(errors, fmt.Errorf("days of week must be %v", validDaysOfWeek))
		return warnings, errors
	}

	return warnings, errors
}

func HoursOfDay(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return warnings, errors
	}

	hoursOfDayValidationMsg := "valid hours of day are: 0 - represents time slot 0:00 - 3:59 UTC - 4 - represents time" +
		"slot 4:00 - 7:59 UTC - 8 - represents time slot 8:00 - 11:59 UTC - 12 - represents time slot" +
		"12:00 - 15:59 UTC - 16 - represents time slot 16:00 - 19:59 UTC - 20 - represents time slot" +
		"20:00 - 23:59 UTC"

	if (v < 0 || v > 20) || (v%4 != 0) {
		errors = append(errors, fmt.Errorf("%s", hoursOfDayValidationMsg))
		return warnings, errors
	}

	return warnings, errors
}

func LeadTimeInWeeks(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return warnings, errors
	}

	if v < 1 || v > 4 {
		errors = append(errors, fmt.Errorf("the lead time in weeks must be between %d and %d", 1, 4))
		return warnings, errors
	}

	return warnings, errors
}

func Month(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	validMonth := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

	if !slices.Contains(validMonth, v) {
		errors = append(errors, fmt.Errorf("month must be %v", validMonth))
		return warnings, errors
	}

	return warnings, errors
}

func WeeksOfMonth(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return warnings, errors
	}

	if v < 1 || v > 4 {
		errors = append(errors, fmt.Errorf("weeksOfMonth should be a list of integers between %d and %d (inclusive)", 1, 4))
		return warnings, errors
	}

	return warnings, errors
}

func Preference(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	if v != string(cloudexadatainfrastructures.PreferenceCustomPreference) && v != string(cloudexadatainfrastructures.PreferenceNoPreference) {
		errors = append(errors, fmt.Errorf("%v must be %v or %v", k,
			string(cloudexadatainfrastructures.PreferenceCustomPreference), string(cloudexadatainfrastructures.PreferenceNoPreference)))
		return warnings, errors
	}

	return warnings, errors
}

func PatchingMode(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	if v != string(cloudexadatainfrastructures.PatchingModeRolling) && v != string(cloudexadatainfrastructures.PatchingModeNonRolling) {
		errors = append(errors, fmt.Errorf("%v must be %v or %v", k,
			string(cloudexadatainfrastructures.PatchingModeRolling), string(cloudexadatainfrastructures.PatchingModeNonRolling)))
		return warnings, errors
	}

	return warnings, errors
}

func ExadataName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	if v == "" {
		errors = append(errors, fmt.Errorf("%v must not be an empty string", k))
		return warnings, errors
	}

	return warnings, errors
}
