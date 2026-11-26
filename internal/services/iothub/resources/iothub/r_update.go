// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id, err := parse.IotHubID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing %s: %+v", id, err)
	}

	locks.ByName(id.Name, IothubResourceName)
	defer locks.UnlockByName(id.Name, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	if iothub.Properties == nil {
		return fmt.Errorf("reading %s: properties was nil", id)
	}

	prop := *iothub.Properties

	if d.HasChange("sku") {
		iothub.Sku = expandIoTHubSku(d)
	}

	if d.HasChange("identity") {
		identity, err := expandIotHubIdentity(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		iothub.Identity = identity
	}

	if d.HasChange("tags") {
		iothub.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("route") {
		if prop.Routing == nil {
			prop.Routing = &devices.RoutingProperties{}
		}
		prop.Routing.Routes = expandIoTHubRoutes(d)
	}

	if d.HasChange("enrichment") {
		if prop.Routing == nil {
			prop.Routing = &devices.RoutingProperties{}
		}
		prop.Routing.Enrichments = expandIoTHubEnrichments(d)
	}

	if d.HasChange("fallback_route") {
		if prop.Routing == nil {
			prop.Routing = &devices.RoutingProperties{}
		}
		if _, ok := d.GetOk("fallback_route"); ok {
			prop.Routing.FallbackRoute = expandIoTHubFallbackRoute(d)
		} else {
			prop.Routing.FallbackRoute = &devices.FallbackRouteProperties{
				Source:        pointer.To(string(devices.RoutingSourceDeviceMessages)),
				Condition:     pointer.To("true"),
				EndpointNames: &[]string{"events"},
				IsEnabled:     pointer.To(true),
			}
		}
	}

	if d.HasChange("endpoint") {
		if prop.Routing == nil {
			prop.Routing = &devices.RoutingProperties{}
		}
		prop.Routing.Endpoints, err = expandIoTHubEndpoints(d, subscriptionId)
		if err != nil {
			return fmt.Errorf("expanding `endpoint`: %+v", err)
		}
	}

	if d.HasChange("file_upload") {
		storageEndpoints, messagingEndpoints, enableFileUploadNotifications, err := expandIoTHubFileUpload(d)
		if err != nil {
			return fmt.Errorf("expanding `file_upload`: %+v", err)
		}
		prop.StorageEndpoints = storageEndpoints
		prop.MessagingEndpoints = messagingEndpoints
		prop.EnableFileUploadNotifications = &enableFileUploadNotifications
	}

	if d.HasChange("cloud_to_device") {
		cloudToDeviceProperties := &devices.CloudToDeviceProperties{}
		if _, ok := d.GetOk("cloud_to_device"); ok {
			cloudToDeviceProperties = expandIoTHubCloudToDevice(d)
		}
		prop.CloudToDevice = cloudToDeviceProperties
	}

	if d.HasChange("network_rule_set") {
		if _, ok := d.GetOk("network_rule_set"); ok {
			prop.NetworkRuleSets = expandNetworkRuleSetProperties(d)
		} else {
			prop.NetworkRuleSets = &devices.NetworkRuleSetProperties{}
		}
	}

	if d.HasChange("public_network_access_enabled") {
		//nolint:staticcheck
		if v, ok := d.GetOkExists("public_network_access_enabled"); ok {
			enabled := devices.PublicNetworkAccessDisabled
			if v.(bool) {
				enabled = devices.PublicNetworkAccessEnabled
			}
			prop.PublicNetworkAccess = enabled
		}
	}

	if d.HasChange("event_hub_retention_in_days") {
		retention, retentionOk := d.GetOk("event_hub_retention_in_days")
		if prop.EventHubEndpoints == nil {
			prop.EventHubEndpoints = make(map[string]*devices.EventHubProperties)
		}
		eh, ok := prop.EventHubEndpoints["events"]
		if !ok {
			prop.EventHubEndpoints["events"] = &devices.EventHubProperties{}
		}
		eh.RetentionTimeInDays = nil
		if retentionOk {
			eh.RetentionTimeInDays = pointer.To(int64(retention.(int)))
		}
	}

	if d.HasChange("event_hub_partition_count") {
		partition, partitionOk := d.GetOk("event_hub_partition_count")
		if prop.EventHubEndpoints == nil {
			prop.EventHubEndpoints = make(map[string]*devices.EventHubProperties)
		}
		if eh, ok := prop.EventHubEndpoints["events"]; ok {
			eh.PartitionCount = nil
			if partitionOk {
				eh.PartitionCount = pointer.To(int32(partition.(int)))
			}
		}
	}

	if d.HasChange("local_authentication_enabled") {
		prop.DisableLocalAuth = pointer.To(!d.Get("local_authentication_enabled").(bool))
	}

	if d.HasChange("min_tls_version") {
		prop.MinTLSVersion = nil
		if v, ok := d.GetOk("min_tls_version"); ok {
			prop.MinTLSVersion = pointer.To(v.(string))
		}
	}

	iothub.Properties = &prop
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %q: %+v", id, err)
	}

	// When the `local_authentication_enabled` updated, GET on the resource is returned with 404 for a while.
	// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/27183
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"404"},
		Target:  []string{"200"},
		Refresh: func() (result any, state string, err error) {
			resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return resp, strconv.Itoa(resp.StatusCode), nil
				}
				return resp, strconv.Itoa(resp.StatusCode), err
			}
			return resp, strconv.Itoa(resp.StatusCode), nil
		},
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
		Timeout:      time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubRead(d, meta)
}
