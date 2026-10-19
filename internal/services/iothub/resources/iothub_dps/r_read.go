// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseProvisioningServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ProvisioningServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if err := d.Set("sku", flattenIoTHubDPSSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		props := model.Properties
		if err := d.Set("linked_hub", flattenIoTHubDPSLinkedHub(props.IotHubs)); err != nil {
			return fmt.Errorf("setting `linked_hub`: %+v", err)
		}

		if err := d.Set("ip_filter_rule", flattenDpsIPFilterRules(props.IPFilterRules)); err != nil {
			return fmt.Errorf("setting `ip_filter_rule` in IoTHub DPS %q: %+v", id.ProvisioningServiceName, err)
		}

		d.Set("service_operations_host_name", props.ServiceOperationsHostName)
		d.Set("device_provisioning_host_name", props.DeviceProvisioningHostName)
		d.Set("id_scope", props.IdScope)

		allocationPolicy := string(iotdpsresource.AllocationPolicyHashed)
		if props.AllocationPolicy != nil {
			allocationPolicy = string(*props.AllocationPolicy)
		}
		d.Set("allocation_policy", allocationPolicy)

		d.Set("data_residency_enabled", pointer.From(props.EnableDataResidency))

		publicNetworkAccess := true
		if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess != "" {
			publicNetworkAccess = strings.EqualFold("Enabled", string(*props.PublicNetworkAccess))
		}
		d.Set("public_network_access_enabled", publicNetworkAccess)

		d.Set("tags", flattenTags(model.Tags))
	}

	return nil
}

func waitForIotHubDPSToBeDeleted(ctx context.Context, client *iotdpsresource.IotDpsResourceClient, id commonids.ProvisioningServiceId) error {
	log.Printf("[DEBUG] Waiting for IoT Device Provisioning Service %q to be deleted", id)
	poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, id)
		return resp.HttpResponse, err
	}, custompollers.DefaultDeletionEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for IoT Device Provisioning Service %q to be deleted: %+v", id, err)
	}
	return nil
}
