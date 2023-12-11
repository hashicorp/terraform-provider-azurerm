// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func IoTHubIpRuleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-:.+%_#*?!(),=@;']{1,128}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("`ip_rule_name` can only be alphanumeric string up to 128 characters long. Only the ASCII 7-bit alphanumeric characters plus the following special characters are accepted: - : . + %% _ # * ? ! ( ) , = @ ; '"))
	}

	return
}
