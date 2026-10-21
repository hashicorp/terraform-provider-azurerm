// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_profile

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkProfileCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkProfiles
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := networkprofiles.NewNetworkProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, networkprofiles.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(azureNetworkProfileResourceName, id.ID())
		}
	}

	containerNetworkInterfaceConfigurations := expandNetworkProfileContainerNetworkInterface(d.Get("container_network_interface").([]any))
	subnetsToLock, vnetsToLock, err := expandNetworkProfileVirtualNetworkSubnetNames(containerNetworkInterfaceConfigurations)
	if err != nil {
		return fmt.Errorf("extracting names of Subnet and Virtual Network: %+v", err)
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	locks.MultipleByID(vnetsToLock)
	defer locks.UnlockMultipleByID(vnetsToLock)

	locks.MultipleByID(subnetsToLock)
	defer locks.UnlockMultipleByID(subnetsToLock)

	payload := networkprofiles.NetworkProfile{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Properties: &networkprofiles.NetworkProfilePropertiesFormat{
			ContainerNetworkInterfaceConfigurations: containerNetworkInterfaceConfigurations,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceNetworkProfileRead(d, meta)
}

func expandNetworkProfileContainerNetworkInterface(input []any) *[]networkprofiles.ContainerNetworkInterfaceConfiguration {
	retCNIConfigs := make([]networkprofiles.ContainerNetworkInterfaceConfiguration, 0)

	for _, cniConfig := range input {
		nciData := cniConfig.(map[string]any)
		ipConfigs := nciData["ip_configuration"].([]any)

		retIPConfigs := make([]networkprofiles.IPConfigurationProfile, 0)
		for _, ipConfig := range ipConfigs {
			ipData := ipConfig.(map[string]any)

			retIPConfig := networkprofiles.IPConfigurationProfile{
				Name: pointer.To(ipData["name"].(string)),
				Properties: &networkprofiles.IPConfigurationProfilePropertiesFormat{
					Subnet: &networkprofiles.Subnet{
						Id: pointer.To(ipData["subnet_id"].(string)),
					},
				},
			}

			retIPConfigs = append(retIPConfigs, retIPConfig)
		}

		retCNIConfig := networkprofiles.ContainerNetworkInterfaceConfiguration{
			Name: pointer.To(nciData["name"].(string)),
			Properties: &networkprofiles.ContainerNetworkInterfaceConfigurationPropertiesFormat{
				IPConfigurations: &retIPConfigs,
			},
		}

		retCNIConfigs = append(retCNIConfigs, retCNIConfig)
	}

	return &retCNIConfigs
}

func expandNetworkProfileVirtualNetworkSubnetNames(input *[]networkprofiles.ContainerNetworkInterfaceConfiguration) (*[]string, *[]string, error) {
	subnetIds := make([]string, 0)
	vnetIds := make([]string, 0)

	if input != nil {
		for _, item := range *input {
			if item.Properties == nil || item.Properties.IPConfigurations == nil {
				continue
			}

			for _, config := range *item.Properties.IPConfigurations {
				if config.Properties == nil || config.Properties.Subnet == nil || config.Properties.Subnet.Id == nil {
					continue
				}

				subnetId, err := commonids.ParseSubnetIDInsensitively(*config.Properties.Subnet.Id)
				if err != nil {
					return nil, nil, err
				}

				vnetId := commonids.NewVirtualNetworkID(subnetId.SubscriptionId, subnetId.ResourceGroupName, subnetId.VirtualNetworkName)

				if !slices.Contains(subnetIds, subnetId.ID()) {
					subnetIds = append(subnetIds, subnetId.ID())
				}

				if !slices.Contains(vnetIds, vnetId.ID()) {
					vnetIds = append(vnetIds, vnetId.ID())
				}
			}
		}
	}

	return &subnetIds, &vnetIds, nil
}
