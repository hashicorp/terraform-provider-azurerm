// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strconv"
	"strings"
)

func ProtocolWithPort(input interface{}, k string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
		return
	}

	parts := strings.Split(v, ":")
	if len(parts) != 2 {
		errors = append(errors, fmt.Errorf("expected %s to be a two part string separated by a `:`, e.g. TCP:80", k))
		return
	}

	if parts[0] != "TCP" && parts[0] != "UDP" {
		errors = append(errors, fmt.Errorf("protocol portion of %s must be one of `TCP` or `UDP`, got %q", k, parts[0]))
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil || port == 0 || port > 65535 {
		errors = append(errors, fmt.Errorf("port in %s must me an integer value between 1 and 65535, got %q", k, parts[1]))
	}

	return
}
