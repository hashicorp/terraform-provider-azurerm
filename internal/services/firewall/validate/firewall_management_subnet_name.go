package validate

import (
	"fmt"

	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
)

func FirewallManagementSubnetName(v interface{}, k string) (warnings []string, errors []error) {
	parsed, err := networkParse.SubnetID(v.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %+v", v.(string), err))
		return warnings, errors
	}

	if parsed.Name != "AzureFirewallManagementSubnet" {
		errors = append(errors, fmt.Errorf("The name of the management subnet for %q must be exactly 'AzureFirewallManagementSubnet' to be used for the Azure Firewall resource", k))
	}

	return warnings, errors
}
