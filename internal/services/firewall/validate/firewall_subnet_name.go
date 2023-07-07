// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

func FirewallSubnetName(v interface{}, k string) (warnings []string, errors []error) {
	parsed, err := commonids.ParseSubnetID(v.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %+v", v.(string), err))
		return warnings, errors
	}

	if parsed.SubnetName != "AzureFirewallSubnet" {
		errors = append(errors, fmt.Errorf("The name of the Subnet for %q must be exactly 'AzureFirewallSubnet' to be used for the Azure Firewall resource", k))
	}

	return warnings, errors
}
