// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_profile

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceNetworkProfileFlatten(d *pluginsdk.ResourceData, id *networkprofiles.NetworkProfileId, profile *networkprofiles.NetworkProfile) error {
	d.Set("name", id.NetworkProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if profile != nil {
		if props := profile.Properties; props != nil {
			if err := d.Set("container_network_interface", flattenNetworkProfileContainerNetworkInterface(props.ContainerNetworkInterfaceConfigurations)); err != nil {
				return fmt.Errorf("setting `container_network_interface`: %+v", err)
			}

			if err := d.Set("container_network_interface_ids", flattenNetworkProfileContainerNetworkInterfaceIDs(props.ContainerNetworkInterfaces)); err != nil {
				return fmt.Errorf("setting `container_network_interface_ids`: %+v", err)
			}
		}
		d.Set("location", location.NormalizeNilable(profile.Location))
		if err := tags.FlattenAndSet(d, profile.Tags); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenNetworkProfileContainerNetworkInterface(input *[]networkprofiles.ContainerNetworkInterfaceConfiguration) []any {
	output := make([]any, 0)
	if input == nil {
		return output
	}

	for _, cniConfig := range *input {
		ipConfigurations := make([]any, 0)
		if props := cniConfig.Properties; props != nil && props.IPConfigurations != nil {
			for _, ipConfig := range *props.IPConfigurations {
				subnetId := ""
				if ipProps := ipConfig.Properties; ipProps != nil && ipProps.Subnet != nil && ipProps.Subnet.Id != nil {
					subnetId = *ipProps.Subnet.Id
				}

				ipConfigurations = append(ipConfigurations, map[string]any{
					"name":      pointer.From(ipConfig.Name),
					"subnet_id": subnetId,
				})
			}
		}

		output = append(output, map[string]any{
			"name":             pointer.From(cniConfig.Name),
			"ip_configuration": ipConfigurations,
		})
	}

	return output
}

func flattenNetworkProfileContainerNetworkInterfaceIDs(input *[]networkprofiles.ContainerNetworkInterface) []string {
	retCNIs := make([]string, 0)

	if input != nil {
		for _, retCNI := range *input {
			if retCNI.Id != nil {
				retCNIs = append(retCNIs, *retCNI.Id)
			}
		}
	}

	return retCNIs
}
