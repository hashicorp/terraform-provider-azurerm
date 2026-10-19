// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_hardware_security_module

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2021-11-30/dedicatedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDedicatedHardwareSecurityModuleCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HSM.DedicatedHsmClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := dedicatedhsms.NewDedicatedHSMID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.DedicatedHsmGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dedicated_hardware_security_module", id.ID())
		}
	}

	skuName := dedicatedhsms.SkuName(d.Get("sku_name").(string))
	if _, ok := d.GetOk("management_network_profile"); ok {
		if skuName == dedicatedhsms.SkuNameSafeNetLunaNetworkHSMASevenNineZero {
			return fmt.Errorf("management_network_profile should not be specified when sku_name is %s", skuName)
		}
	}

	parameters := dedicatedhsms.DedicatedHsm{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: dedicatedhsms.DedicatedHsmProperties{
			NetworkProfile:           expandDedicatedHsmNetworkProfile(d.Get("network_profile").([]any)),
			ManagementNetworkProfile: expandDedicatedHsmNetworkProfile(d.Get("management_network_profile").([]any)),
		},
		Sku: &dedicatedhsms.Sku{
			Name: &skuName,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("stamp_id"); ok {
		parameters.Properties.StampId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		zones := zones.ExpandUntyped(v.(*pluginsdk.Set).List())
		if len(zones) > 0 {
			parameters.Zones = &zones
		}
	}

	if err := client.DedicatedHsmCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDedicatedHardwareSecurityModuleRead(d, meta)
}

func expandDedicatedHsmNetworkProfile(input []any) *dedicatedhsms.NetworkProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	result := dedicatedhsms.NetworkProfile{
		Subnet: &dedicatedhsms.ApiEntityReference{
			Id: pointer.To(v["subnet_id"].(string)),
		},
		NetworkInterfaces: expandDedicatedHsmNetworkInterfacePrivateIPAddresses(v["network_interface_private_ip_addresses"].(*pluginsdk.Set).List()),
	}

	return &result
}

func expandDedicatedHsmNetworkInterfacePrivateIPAddresses(input []any) *[]dedicatedhsms.NetworkInterface {
	results := make([]dedicatedhsms.NetworkInterface, 0)

	for _, item := range input {
		if item != nil {
			results = append(results, dedicatedhsms.NetworkInterface{
				PrivateIPAddress: pointer.To(item.(string)),
			})
		}
	}

	return &results
}
