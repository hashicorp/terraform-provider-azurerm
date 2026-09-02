// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// MaintenanceWindow validation

func DaysOfWeek(i interface{}, k string) ([]string, []error) {
	return validation.IsDayOfTheWeek(false)(i, k)
}

// lintignore:V012 // valid values are multiples of 4, error message documents the time slots
// intentionally not validation.IntInSlice([0,4,8,12,16,20]): its generic error would lose the
// message below documenting which maintenance-window time slot each value represents
func HoursOfDay(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	hoursOfDayValidationMsg := "valid hours of day are: 0 - represents time slot 0:00 - 3:59 UTC - 4 - represents time" +
		"slot 4:00 - 7:59 UTC - 8 - represents time slot 8:00 - 11:59 UTC - 12 - represents time slot" +
		"12:00 - 15:59 UTC - 16 - represents time slot 16:00 - 19:59 UTC - 20 - represents time slot" +
		"20:00 - 23:59 UTC"

	if (v < 0 || v > 20) || (v%4 != 0) {
		errors = append(errors, fmt.Errorf("%s", hoursOfDayValidationMsg))
		return
	}

	return
}

func Month(i interface{}, k string) ([]string, []error) {
	return validation.IsMonth(false)(i, k)
}

func Preference(i interface{}, k string) ([]string, []error) {
	return validation.StringInSlice([]string{
		string(cloudexadatainfrastructures.PreferenceCustomPreference),
		string(cloudexadatainfrastructures.PreferenceNoPreference),
	}, false)(i, k)
}

func PatchingMode(i interface{}, k string) ([]string, []error) {
	return validation.StringInSlice([]string{
		string(cloudexadatainfrastructures.PatchingModeRolling),
		string(cloudexadatainfrastructures.PatchingModeNonRolling),
	}, false)(i, k)
}
