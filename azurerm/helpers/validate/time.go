package validate

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	iso8601 "github.com/btubbs/datetime"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rickb777/date/period"
)

func ISO8601Duration(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if _, err := period.Parse(v); err != nil {
		errors = append(errors, err)
	}
	return warnings, errors
}

func ISO8601DurationBetween(min string, max string) func(i interface{}, k string) (warnings []string, errors []error) {
	minDuration := period.MustParse(min).DurationApprox()
	maxDuration := period.MustParse(max).DurationApprox()
	if minDuration >= maxDuration {
		panic(fmt.Sprintf("min duration (%v) >= max duration (%v)", minDuration, maxDuration))
	}
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
		}

		p, err := period.Parse(v)
		if err != nil {
			return nil, []error{err}
		}

		duration := p.DurationApprox()
		if duration < minDuration || duration > maxDuration {
			return nil, []error{fmt.Errorf("expected %s to be in the range (%v - %v), got %v", k, minDuration, maxDuration, duration)}
		}

		return nil, nil
	}
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
