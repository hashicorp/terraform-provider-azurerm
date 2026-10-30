// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewProvisioningServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	iotdps := resp.Model
	if iotdps == nil {
		return fmt.Errorf("retrieving model of %s: %+v", id, err)
	}

	if d.HasChanges("allocation_policy") {
		allocationPolicy := iotdpsresource.AllocationPolicy(d.Get("allocation_policy").(string))
		iotdps.Properties.AllocationPolicy = &allocationPolicy
	}

	if d.HasChanges("public_network_access_enabled") {
		publicNetworkAccess := iotdpsresource.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = iotdpsresource.PublicNetworkAccessDisabled
		}
		iotdps.Properties.PublicNetworkAccess = &publicNetworkAccess
	}

	if d.HasChanges("ip_filter_rule") {
		iotdps.Properties.IPFilterRules = expandDpsIPFilterRules(d)
	}

	if d.HasChanges("linked_hub") {
		iotdps.Properties.IotHubs = expandIoTHubDPSIoTHubs(d.Get("linked_hub").([]any))
	}

	if d.HasChanges("sku") {
		iotdps.Sku = expandIoTHubDPSSku(d)
	}

	if d.HasChanges("tags") {
		iotdps.Tags = expandTags(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, *iotdps); err != nil {
		return fmt.Errorf("updating IoT Device Provisioning Service %s: %+v", id, err)
	}

	return resourceIotHubDPSRead(d, meta)
}
