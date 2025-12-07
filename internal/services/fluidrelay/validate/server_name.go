// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

var serverNameReg = regexp.MustCompile(`^[-0-9a-zA-Z]{1,50}$`)

func FluidRelayServerName(input interface{}, key string) (warnings []string, errs []error) {
	v, ok := input.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected %q to be a string", key))
		return
	}
	// Name should contain only alphanumeric characters and hyphens, up to 50 characters long.
	if !serverNameReg.MatchString(v) {
		errs = append(errs, fmt.Errorf("%s should contain only alphanumeric characters and hyphens, up to 50 characters long", key))
	}
	return
}
