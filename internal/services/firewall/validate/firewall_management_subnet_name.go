// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

func FirewallManagementSubnetName(v interface{}, k string) (warnings []string, errors []error) {
	parsed, err := commonids.ParseSubnetID(v.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %+v", v.(string), err))
		return warnings, errors
	}

	if !strings.EqualFold(parsed.SubnetName, "AzureFirewallManagementSubnet") {
		errors = append(errors, fmt.Errorf("the name of the management subnet for %q must be 'AzureFirewallManagementSubnet' (case-insensitive) to be used for the Azure Firewall resource", k))
	}

	return warnings, errors
}
