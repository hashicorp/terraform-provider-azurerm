// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_hardware_security_module

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2021-11-30/dedicatedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHardwareSecurityModuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dedicatedhsms.ParseDedicatedHSMID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DedicatedHsmGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DedicatedHSMName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		props := model.Properties

		if err := d.Set("management_network_profile", flattenDedicatedHsmNetworkProfile(props.ManagementNetworkProfile)); err != nil {
			return fmt.Errorf("setting management_network_profile: %+v", err)
		}

		if err := d.Set("network_profile", flattenDedicatedHsmNetworkProfile(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting network_profile: %+v", err)
		}
		d.Set("stamp_id", props.StampId)

		skuName := ""
		if model.Sku.Name != nil {
			skuName = string(*model.Sku.Name)
		}
		d.Set("sku_name", skuName)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenDedicatedHsmNetworkProfile(input *dedicatedhsms.NetworkProfile) []any {
	if input == nil {
		return make([]any, 0)
	}

	var subnetId string
	if input.Subnet != nil && input.Subnet.Id != nil {
		subnetId = *input.Subnet.Id
	}

	return []any{
		map[string]any{
			"network_interface_private_ip_addresses": flattenDedicatedHsmNetworkInterfacePrivateIPAddresses(input.NetworkInterfaces),
			"subnet_id":                              subnetId,
		},
	}
}

func flattenDedicatedHsmNetworkInterfacePrivateIPAddresses(input *[]dedicatedhsms.NetworkInterface) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.PrivateIPAddress != nil {
			results = append(results, *item.PrivateIPAddress)
		}
	}

	return results
}
