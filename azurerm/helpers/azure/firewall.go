package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
)

// The API requires InternalPublicIPAddress to be set when for a CreateOrUpdate
// operation, but Get operations return the property as PublicIPAddress
// so we need to go through and copy the value to the correct property.
func FirewallFixIPConfiguration(input *[]network.AzureFirewallIPConfiguration) (*[]network.AzureFirewallIPConfiguration, error) {
	if input == nil {
		return nil, fmt.Errorf("`input` was nil")
	}

	results := make([]network.AzureFirewallIPConfiguration, 0)
	for _, config := range *input {
		if config.Subnet == nil || config.Subnet.ID == nil {
			return nil, fmt.Errorf("`config.Subnet.ID` was nil")
		}

		if config.PublicIPAddress == nil || config.PublicIPAddress.ID == nil {
			return nil, fmt.Errorf("`config.PublicIPAddress.ID` was nil")
		}

		result := network.AzureFirewallIPConfiguration{
			Name: config.Name,
			AzureFirewallIPConfigurationPropertiesFormat: &network.AzureFirewallIPConfigurationPropertiesFormat{
				Subnet: &network.SubResource{
					ID: config.Subnet.ID,
				},
				InternalPublicIPAddress: &network.SubResource{
					ID: config.PublicIPAddress.ID,
				},
			},
		}
		results = append(results, result)
	}

	return &results, nil
}
