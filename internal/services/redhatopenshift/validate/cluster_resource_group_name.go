// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func ClusterResourceGroupName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) > 90 {
		errors = append(errors, fmt.Errorf("%q may not exceed 90 characters in length", k))
	}

	if strings.HasSuffix(value, ".") {
		errors = append(errors, fmt.Errorf("%q may not end with a period", k))
	}

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be blank", k))
	} else if matched := regexp.MustCompile(`^[0-9a-z-._()]+$`).Match([]byte(value)); !matched {
		// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
		// ARO only allow for lower cases https://github.com/Azure/ARO-RP/blob/e5c40654277c77fe78ba669610ac05774e448683/pkg/frontend/openshiftcluster_putorpatch.go#L189
		errors = append(errors, fmt.Errorf("%q may only contain lowercase alpha characters, digit, dash, underscores, parentheses and periods", k))
	}

	return warnings, errors
}
