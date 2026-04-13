// Copyright IBM Corp. 2014, 2025
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

	if v == "any" || v == "application-default" {
		return
	}

	parts := strings.Split(v, ":")
	if len(parts) != 2 {
		errors = append(errors, fmt.Errorf("expected %s to be a two part string separated by a `:`, e.g. TCP:80, or a supported keyword like `any`", k))
		return
	}

	if parts[0] != "TCP" && parts[0] != "UDP" {
		errors = append(errors, fmt.Errorf("protocol portion of %s must be one of `TCP` or `UDP`, got %q", k, parts[0]))
	}

	if strings.Contains(parts[1], "-") {
		rangeParts := strings.Split(parts[1], "-")
		if len(rangeParts) != 2 {
			errors = append(errors, fmt.Errorf("port range in %s must be in format START-END, e.g. TCP:1024-1206", k))
			return
		}
		startPort, err := strconv.Atoi(rangeParts[0])
		if err != nil || startPort < 1 || startPort > 65535 {
			errors = append(errors, fmt.Errorf("start port in %s must be an integer between 1 and 65535, got %q", k, rangeParts[0]))
			return
		}
		endPort, err := strconv.Atoi(rangeParts[1])
		if err != nil || endPort < 1 || endPort > 65535 {
			errors = append(errors, fmt.Errorf("end port in %s must be an integer between 1 and 65535, got %q", k, rangeParts[1]))
			return
		}
		if startPort > endPort {
			errors = append(errors, fmt.Errorf("start port must be less than or equal to end port in %s", k))
			return
		}
		return
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil || port == 0 || port > 65535 {
		errors = append(errors, fmt.Errorf("port in %s must be an integer value between 1 and 65535, got %q", k, parts[1]))
	}

	return
}
