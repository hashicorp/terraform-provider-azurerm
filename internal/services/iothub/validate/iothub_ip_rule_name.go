// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"errors"
	"regexp"
)

func IoTHubIpRuleName(v any, k string) (warnings []string, errs []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z-:.+%_#*?!(),=@;']{1,128}$`).Match([]byte(value)); !matched {
		errs = append(errs, errors.New("`ip_rule_name` can only be alphanumeric string up to 128 characters long. Only the ASCII 7-bit alphanumeric characters plus the following special characters are accepted: - : . + %% _ # * ? ! ( ) , = @ ; '"))
	}

	return
}
