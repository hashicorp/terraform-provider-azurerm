// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func StatusCodeRange(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}
	parts := strings.Split(v, "-")
	if len(parts) < 1 || len(parts) > 2 {
		errors = append(errors, fmt.Errorf("%q must be either a single HTTP status code or a range in the form 100-599", k))
	}

	for _, part := range parts {
		if matched := regexp.MustCompile(`^([1-5][0-9][0-9])$`).Match([]byte(part)); !matched {
			errors = append(errors, fmt.Errorf("%q must be either a single HTTP status code or a range in the form 100-599", k))
		}
	}

	if len(parts) == 2 {
		lowCode, err := strconv.Atoi(parts[0])
		if err != nil {
			errors = append(errors, fmt.Errorf("could not convert status code low value (%+v) to int", parts[1]))
		}
		highCode, err := strconv.Atoi(parts[1])
		if err != nil {
			errors = append(errors, fmt.Errorf("could not convert status code high value (%+v) to int", parts[1]))
		}
		if lowCode > highCode {
			errors = append(errors, fmt.Errorf("%q range values must be in the form low to high, e.g. 200-30. Got %+v", k, v))
		}
	}

	return
}
