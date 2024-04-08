// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func ConfigurationStoreReplicaName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[a-zA-Z0-9]{1,50}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q Replica name may only contain alphanumeric characters and must be between 1-50 chars", k))
	}

	return warnings, errors
}
