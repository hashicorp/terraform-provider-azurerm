// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"errors"
	"fmt"
	"regexp"
)

func SubscriptionName(i interface{}, k string) (warnings []string, errs []error) {
	v, ok := i.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if len(v) > 64 || v == "" {
		errs = append(errs, errors.New("Subscription Name must be between 1 and 64 characters in length"))
	}

	if regexp.MustCompile("[<>;|]").MatchString(v) {
		errs = append(errs, errors.New("Subsciption Name cannot contain the characters `<`, `>`, `;`, or `|`"))
	}

	return
}
