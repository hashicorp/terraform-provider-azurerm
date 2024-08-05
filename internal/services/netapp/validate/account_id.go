// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
)

func ValidateNetAppAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := netappaccounts.ParseNetAppAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
