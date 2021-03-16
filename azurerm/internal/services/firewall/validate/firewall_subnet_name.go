package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func FirewallSubnetName(v interface{}, k string) (warnings []string, errors []error) {
	parsed, err := azure.ParseAzureResourceID(v.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("Error parsing Azure Resource ID %q", v.(string)))
		return warnings, errors
	}
	subnetName := parsed.Path["subnets"]
	if subnetName != "AzureFirewallSubnet" {
		errors = append(errors, fmt.Errorf("The name of the Subnet for %q must be exactly 'AzureFirewallSubnet' to be used for the Azure Firewall resource", k))
	}

	return warnings, errors
}
