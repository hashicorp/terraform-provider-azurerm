// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_network_interface

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.NetworkInterfaces

			id, err := networkinterfaces.ParseNetworkInterfaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schema := StackHCINetworkInterfaceResourceModel{
				Name:              id.NetworkInterfaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				if model.ExtendedLocation != nil && model.ExtendedLocation.Name != nil {
					customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(*model.ExtendedLocation.Name)
					if err != nil {
						return err
					}

					schema.CustomLocationId = customLocationId.ID()
				}

				if props := model.Properties; props != nil {
					schema.MACAddress = pointer.From(props.MacAddress)

					ipConfiguration, err := flattenStackHCINetworkInterfaceIPConfiguration(props.IPConfigurations)
					if err != nil {
						return err
					}
					schema.IPConfiguration = ipConfiguration

					if props.DnsSettings != nil {
						schema.DNSServers = pointer.From(props.DnsSettings.DnsServers)
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func flattenStackHCINetworkInterfaceIPConfiguration(input *[]networkinterfaces.IPConfiguration) ([]StackHCIIPConfigurationModel, error) {
	if input == nil {
		return make([]StackHCIIPConfigurationModel, 0), nil
	}

	results := make([]StackHCIIPConfigurationModel, 0)
	for _, v := range *input {
		result := StackHCIIPConfigurationModel{}

		if v.Properties != nil {
			result.Gateway = pointer.From(v.Properties.Gateway)
			result.PrefixLength = pointer.From(v.Properties.PrefixLength)
			result.PrivateIPAddress = pointer.From(v.Properties.PrivateIPAddress)

			var subnetId string
			if v.Properties.Subnet != nil && v.Properties.Subnet.Id != nil {
				parsedSubnetId, err := logicalnetworks.ParseLogicalNetworkIDInsensitively(*v.Properties.Subnet.Id)
				if err != nil {
					return results, err
				}

				subnetId = parsedSubnetId.ID()
			}
			result.SubnetID = subnetId

			results = append(results, result)
		}
	}

	return results, nil
}
