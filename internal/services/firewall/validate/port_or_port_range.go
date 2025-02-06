// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strconv"
)

func PortOrPortRangeWithin(min int, max int) func(interface{}, string) ([]string, []error) {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
			return warnings, errors
		}

		assertWithinRange := func(n int) error {
			if n < min || n > max {
				return fmt.Errorf("port %d is out of range (%d-%d)", n, min, max)
			}

			return nil
		}

		// Allowed format including: `num` or `num1-num2` (num1 < num2).
		groups := regexp.MustCompile(`^(\d+)((-)(\d+))?$`).FindStringSubmatch(v)
		if len(groups) != 5 {
			errors = append(errors, fmt.Errorf("expected `number` or `num1-num2` but got %q", v))
			return warnings, errors
		}

		if groups[2] == "" {
			p1, _ := strconv.Atoi(groups[1])

			if err := assertWithinRange(p1); err != nil {
				errors = append(errors, err)
				return warnings, errors
			}
		} else {
			p1, _ := strconv.Atoi(groups[1])
			p2, _ := strconv.Atoi(groups[4])

			if p1 >= p2 {
				errors = append(errors, fmt.Errorf("beginning port (%d) should be less than ending port (%d)", p1, p2))
				return warnings, errors
			}

			if err := assertWithinRange(p1); err != nil {
				errors = append(errors, err)
				return warnings, errors
			}

			if err := assertWithinRange(p2); err != nil {
				errors = append(errors, err)
				return warnings, errors
			}
		}

		return nil, nil
	}
}
