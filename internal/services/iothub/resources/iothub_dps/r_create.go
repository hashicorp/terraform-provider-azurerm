// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewProvisioningServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing IoT Device Provisioning Service %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_iothub_dps", id.ID())
		}
	}

	publicNetworkAccess := iotdpsresource.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = iotdpsresource.PublicNetworkAccessDisabled
	}

	allocationPolicy := iotdpsresource.AllocationPolicy(d.Get("allocation_policy").(string))
	iotdps := iotdpsresource.ProvisioningServiceDescription{
		Location: location.Normalize(d.Get("location").(string)),
		Name:     pointer.To(id.ProvisioningServiceName),
		Sku:      expandIoTHubDPSSku(d),
		Properties: iotdpsresource.IotDpsPropertiesDescription{
			IotHubs:             expandIoTHubDPSIoTHubs(d.Get("linked_hub").([]any)),
			AllocationPolicy:    &allocationPolicy,
			EnableDataResidency: pointer.To(d.Get("data_residency_enabled").(bool)),
			IPFilterRules:       expandDpsIPFilterRules(d),
			PublicNetworkAccess: &publicNetworkAccess,
		},
		Tags: expandTags(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, iotdps, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating IoT Device Provisioning Service %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubDPSRead(d, meta)
}
