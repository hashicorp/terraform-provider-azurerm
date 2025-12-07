// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"time"
)

// EndDateTime validates the NGINXaaS Dataplane API Key end_date_time value.
// This value needs to be:
// - a correctly formatted RFC3339 date-time
// - not expired
// - no further out than 2 years
func EndDateTime(i interface{}, k string) (warnings []string, errors []error) {
	var err error
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	var expTime time.Time
	if expTime, err = time.Parse(time.RFC3339, v); err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a valid RFC3339 date, got %q: %+v", k, i, err))
		return warnings, errors
	}

	if expTime.Before(time.Now()) {
		errors = append(errors, fmt.Errorf("expected %q to be a valid RFC3339 date that has not already passed", k))
		return warnings, errors
	}

	if expTime.After(time.Now().AddDate(2, 0, 0)) {
		errors = append(errors, fmt.Errorf("expected %q to be a valid RFC3339 date that is no more than 2 years from now", k))
		return warnings, errors
	}

	return nil, nil
}
