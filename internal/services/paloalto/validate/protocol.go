// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strconv"
	"strings"
)

// protocolWithPortStandaloneValues are values that `ProtocolWithPort` accepts on their
// own, without a `TCP:`/`UDP:` prefix.
var protocolWithPortStandaloneValues = []string{"any", "application-default"}

func ProtocolWithPort(input interface{}, k string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
		return
	}

	for _, standalone := range protocolWithPortStandaloneValues {
		if v == standalone {
			return
		}
	}

	parts := strings.SplitN(v, ":", 2)
	if len(parts) != 2 {
		errors = append(errors, fmt.Errorf("expected %s to be `any`, `application-default`, or a two part string separated by a `:`, e.g. `TCP:80` or `TCP:1024-1206`, got %q", k, v))
		return
	}

	if parts[0] != "TCP" && parts[0] != "UDP" {
		errors = append(errors, fmt.Errorf("protocol portion of %s must be one of `TCP` or `UDP`, got %q", k, parts[0]))
	}

	if err := validatePortOrPortRange(parts[1]); err != nil {
		errors = append(errors, fmt.Errorf("port portion of %s is invalid: %s", k, err))
	}

	return
}

// validatePortOrPortRange validates that input is either a single port (`80`) or an
// increasing range of ports (`1024-1206`), with each bound between 1 and 65535.
func validatePortOrPortRange(input string) error {
	bounds := strings.SplitN(input, "-", 2)

	start, err := parsePort(bounds[0])
	if err != nil {
		return err
	}

	if len(bounds) == 1 {
		return nil
	}

	end, err := parsePort(bounds[1])
	if err != nil {
		return err
	}

	if start > end {
		return fmt.Errorf("range %q must start at or before its end", input)
	}

	return nil
}

func parsePort(input string) (int, error) {
	port, err := strconv.Atoi(input)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("must be an integer value between 1 and 65535, got %q", input)
	}
	return port, nil
}
