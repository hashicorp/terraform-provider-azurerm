// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strconv"
	"strings"
)

func IpTrafficPort(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if value == "*" {
		return
	}

	err := fmt.Errorf("%q must be a single port between `0` and `65535` or a range in the format `start-end`, or wildcard(`*`)", k)

	if strings.Contains(value, "-") {
		ports := strings.Split(value, "-")
		if len(ports) != 2 {
			errors = append(errors, err)
			return
		}
		if !isValidPortRange(ports[0], ports[1]) {
			errors = append(errors, err)
			return
		}
		return
	}

	if !isValidPort(value) {
		errors = append(errors, err)
		return
	}

	return warnings, errors
}

func isValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false
	}
	return port >= 0 && port <= 65535
}

func isValidPortRange(startStr, endStr string) bool {
	start, err1 := strconv.Atoi(startStr)
	end, err2 := strconv.Atoi(endStr)
	if err1 != nil || err2 != nil {
		return false
	}
	return start >= 0 && start <= 65535 && end >= 0 && end <= 65535 && start < end
}
