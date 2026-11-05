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
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.NetworkInterfaces

			var config StackHCINetworkInterfaceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := networkinterfaces.NewNetworkInterfaceID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := networkinterfaces.NetworkInterfaces{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				ExtendedLocation: &networkinterfaces.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(networkinterfaces.ExtendedLocationTypesCustomLocation),
				},
				Properties: &networkinterfaces.NetworkInterfaceProperties{
					IPConfigurations: expandStackHCINetworkInterfaceIPConfiguration(config.IPConfiguration),
				},
			}

			if config.MACAddress != "" {
				payload.Properties.MacAddress = pointer.To(config.MACAddress)
			}

			if len(config.DNSServers) != 0 {
				payload.Properties.DnsSettings = &networkinterfaces.InterfaceDNSSettings{
					DnsServers: pointer.To(config.DNSServers),
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("performing create %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandStackHCINetworkInterfaceIPConfiguration(input []StackHCIIPConfigurationModel) *[]networkinterfaces.IPConfiguration {
	if len(input) == 0 {
		return nil
	}

	results := make([]networkinterfaces.IPConfiguration, 0)
	for _, v := range input {
		result := networkinterfaces.IPConfiguration{
			Properties: &networkinterfaces.IPConfigurationProperties{
				Subnet: &networkinterfaces.IPConfigurationPropertiesSubnet{
					Id: pointer.To(v.SubnetID),
				},
			},
		}

		if v.PrivateIPAddress != "" {
			result.Properties.PrivateIPAddress = pointer.To(v.PrivateIPAddress)
		}

		results = append(results, result)
	}

	return &results
}
