package validate

import (
	"fmt"

	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
)

func BastionSubnetName(v interface{}, k string) (warnings []string, errors []error) {
	parsed, err := networkParse.SubnetID(v.(string))
	if err != nil {
		errors = append(errors, fmt.Errorf("parsing %q: %+v", v.(string), err))
		return warnings, errors
	}

	if parsed.Name != "AzureBastionSubnet" {
		errors = append(errors, fmt.Errorf("The name of the Subnet for %q must be exactly 'AzureBastionSubnet' to be used for the Azure Bastion Host resource", k))
	}

	return warnings, errors
}
