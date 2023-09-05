// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumes"
)

func VolumeID(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}

	if _, err := volumes.ParseVolumeID(value); err != nil {
		errors = append(errors, err)
	}

	return warnings, errors
}
