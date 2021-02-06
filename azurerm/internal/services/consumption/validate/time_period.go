package validate

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ConsumptionBudgetTimePeriodStartDate(i interface{}, k string) (warnings []string, errors []error) {
	validateRFC3339TimeWarnings, validateRFC3339TimeErrors := validation.IsRFC3339Time(i, k)
	errors = append(errors, validateRFC3339TimeErrors...)
	warnings = append(warnings, validateRFC3339TimeWarnings...)

	if len(errors) != 0 || len(warnings) != 0 {
		return warnings, errors
	}

	// Errors were already checked by validation.IsRFC3339Time
	startDate, _ := date.ParseTime(time.RFC3339, i.(string))

	// The start date must be first of the month
	if startDate.Day() != 1 {
		errors = append(errors, fmt.Errorf("%q must be first of the month, got day %d", k, startDate.Day()))
		return warnings, errors
	}

	// Budget start date must be on or after June 1, 2017.
	earliestPossibleStartDateString := "2017-06-01T00:00:00Z"
	earliestPossibleStartDate, _ := date.ParseTime(time.RFC3339, earliestPossibleStartDateString)
	if startDate.Before(earliestPossibleStartDate) {
		errors = append(errors, fmt.Errorf("%q must be on or after June 1, 2017, got %q", k, i.(string)))
		return warnings, errors
	}

	// Future start date should not be more than twelve months.
	if startDate.After(time.Now().AddDate(0, 12, 0)) {
		warnings = append(warnings, fmt.Sprintf("%q should not be more than twelve months in the future", k))
	}

	return warnings, errors
}
