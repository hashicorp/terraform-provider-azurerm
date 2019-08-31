package validate

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	iso8601 "github.com/btubbs/datetime"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func ISO8601Duration(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	matched, _ := regexp.MatchString(`^P([0-9]+Y)?([0-9]+M)?([0-9]+W)?([0-9]+D)?(T([0-9]+H)?([0-9]+M)?([0-9]+(\.?[0-9]+)?S)?)?$`, v)

	if !matched {
		errors = append(errors, fmt.Errorf("expected %s to be in ISO 8601 duration format, got %s", k, v))
	}
	return warnings, errors
}

//todo, now in terraform helper, switch over once vended
// -> https://github.com/hashicorp/terraform/blob/master/helper/validation/validation.go#L263
func RFC3339Time(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := date.ParseTime(time.RFC3339, v); err != nil {
		errors = append(errors, fmt.Errorf("%q has the invalid RFC3339 date format %q: %+v", k, i, err))
	}

	return warnings, errors
}

func ISO8601DateTime(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := iso8601.Parse(v, time.UTC); err != nil {
		errors = append(errors, fmt.Errorf("%q has the invalid ISO8601 date format %q: %+v", k, i, err))
	}

	return warnings, errors
}

// RFC3339 date is duration d or greater into the future
func RFC3339DateInFutureBy(d time.Duration) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
			return
		}

		t, err := date.ParseTime(time.RFC3339, v)
		if err != nil {
			errors = append(errors, fmt.Errorf("%q has the invalid RFC3339 date format %q: %+v", k, i, err))
			return
		}

		if time.Until(t) < d {
			errors = append(errors, fmt.Errorf("%q is %q and should be at least %q in the future", k, i, d))
		}

		return warnings, errors
	}
}

func DayOfTheWeek(ignoreCase bool) schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	}, ignoreCase)
}

func Month(ignoreCase bool) schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}, ignoreCase)
}
