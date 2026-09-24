// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"fmt"
	"strings"
	"time"

	iso8601 "github.com/btubbs/datetime"
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

func ISO8601RepeatingTime(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if !strings.HasPrefix(v, "R/") {
		errors = append(errors, fmt.Errorf("%s must start with 'R/'", k))
		return
	}

	partsWithoutPrefix := strings.TrimPrefix(v, "R/")

	pIndex := strings.Index(partsWithoutPrefix, "/P")
	if pIndex == -1 {
		errors = append(errors, fmt.Errorf("%s must end with duration", k))
		return
	}

	dateTime := partsWithoutPrefix[:pIndex]

	if _, err := iso8601.Parse(dateTime, time.UTC); err != nil {
		errors = append(errors, fmt.Errorf("%q has the invalid ISO8601 date format %q: %+v", k, i, err))
		return
	}

	if _, err := period.Parse(partsWithoutPrefix[pIndex+1:]); err != nil {
		errors = append(errors, err)
		return
	}

	return warnings, errors
}
