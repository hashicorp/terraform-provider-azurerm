// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

func BastionSubnetName(v interface{}, k string) (warnings []string, errors []error) {
	parsed, err := commonids.ParseSubnetID(v.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %+v", v.(string), err))
		return warnings, errors
	}

	if parsed.SubnetName != "AzureBastionSubnet" {
		errors = append(errors, fmt.Errorf("The name of the Subnet for %q must be exactly 'AzureBastionSubnet' to be used for the Azure Bastion Host resource", k))
	}

	return warnings, errors
}
