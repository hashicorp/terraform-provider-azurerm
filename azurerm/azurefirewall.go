package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
)

// The API requires InternalPublicIPAddress to be set when for a CreateOrUpdate
// operation, but Get operations return the property as PublicIPAddress
// so we need to go through and copy the value to the correct property.
func fixArmAzureFirewallIPConfiguration(firewall *network.AzureFirewall) []network.AzureFirewallIPConfiguration {
	current := *firewall.IPConfigurations
	ipConfigs := make([]network.AzureFirewallIPConfiguration, 0, len(current))
	for _, config := range current {
		properties := network.AzureFirewallIPConfigurationPropertiesFormat{
			Subnet: &network.SubResource{
				ID: config.Subnet.ID,
			},
			InternalPublicIPAddress: &network.SubResource{
				ID: config.PublicIPAddress.ID,
			},
		}
		ipConfig := network.AzureFirewallIPConfiguration{
			Name: config.Name,
			AzureFirewallIPConfigurationPropertiesFormat: &properties,
		}
		ipConfigs = append(ipConfigs, ipConfig)
	}
	return ipConfigs
}

func expandArmAzureFirewallSet(r *schema.Set) *[]string {
	var result []string
	for _, v := range r.List() {
		s := v.(string)
		result = append(result, s)
	}
	return &result
}
