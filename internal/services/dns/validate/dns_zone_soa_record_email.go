// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// lintignore:V013,V011,V001 // false positive - this validates an email; the string comparisons check for empty segments
func DnsZoneSOARecordEmail(v any, k string) (warnings []string, errs []error) {
	value := v.(string)

	if len(value) == 0 {
		errs = append(errs, fmt.Errorf("%q cannot be an empty string: %q", k, v))
		return warnings, errs
	}

	vSegments := strings.Split(strings.TrimSuffix(value, "."), ".")
	if len(vSegments) < 2 || len(vSegments) > 34 {
		errs = append(errs, fmt.Errorf("%q must be between 2 and 34 segments", k))
		return warnings, errs
	}

	for _, segment := range vSegments {
		if segment == "" {
			errs = append(errs, fmt.Errorf("%q cannot contain consecutive period", k))
			return warnings, errs
		}

		if len(segment) > 63 {
			errs = append(errs, errors.New("the each segment of the `email` must contain between 1 and 63 characters"))
			return warnings, errs
		}
	}

	if !regexp.MustCompile(`^[a-zA-Z\d._-]+$`).MatchString(value) {
		errs = append(errs, fmt.Errorf("%q only contains letters, numbers, underscores, dashes and periods", k))
		return warnings, errs
	}

	return warnings, errs
}
