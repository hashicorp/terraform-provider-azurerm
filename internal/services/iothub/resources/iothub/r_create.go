// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	id := parse.NewIotHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.Name, IothubResourceName)
	defer locks.UnlockByName(id.Name, IothubResourceName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_iothub", id.ID())
		}
	}

	res, err := client.CheckNameAvailability(ctx, devices.OperationInputs{Name: &id.Name})
	if err != nil {
		return fmt.Errorf("an error occurred checking if the IoTHub name was unique: %+v", err)
	}

	if !*res.NameAvailable {
		if _, err = client.Get(ctx, id.ResourceGroup, id.Name); err == nil {
			return fmt.Errorf("an IoTHub already exists with the name %q - please choose an alternate name: %s", id.Name, string(res.Reason))
		}
	}

	routingProperties := devices.RoutingProperties{}

	if _, ok := d.GetOk("route"); ok {
		routingProperties.Routes = expandIoTHubRoutes(d)
	}

	if _, ok := d.GetOk("enrichment"); ok {
		routingProperties.Enrichments = expandIoTHubEnrichments(d)
	}

	if _, ok := d.GetOk("fallback_route"); ok {
		routingProperties.FallbackRoute = expandIoTHubFallbackRoute(d)
	} else {
		routingProperties.FallbackRoute = &devices.FallbackRouteProperties{
			Source:        pointer.To(string(devices.RoutingSourceDeviceMessages)),
			Condition:     pointer.To("true"),
			EndpointNames: &[]string{"events"},
			IsEnabled:     pointer.To(true),
		}
	}

	if _, ok := d.GetOk("endpoint"); ok {
		var err error
		routingProperties.Endpoints, err = expandIoTHubEndpoints(d, subscriptionId)
		if err != nil {
			return fmt.Errorf("expanding `endpoint`: %+v", err)
		}
	}

	storageEndpoints, messagingEndpoints, enableFileUploadNotifications, err := expandIoTHubFileUpload(d)
	if err != nil {
		return fmt.Errorf("expanding `file_upload`: %+v", err)
	}

	cloudToDeviceProperties := &devices.CloudToDeviceProperties{}
	if _, ok := d.GetOk("cloud_to_device"); ok {
		cloudToDeviceProperties = expandIoTHubCloudToDevice(d)
	}

	identity, err := expandIotHubIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	props := devices.IotHubDescription{
		Name:     pointer.To(id.Name),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Sku:      expandIoTHubSku(d),
		Properties: &devices.IotHubProperties{
			Routing:                       &routingProperties,
			StorageEndpoints:              storageEndpoints,
			MessagingEndpoints:            messagingEndpoints,
			EnableFileUploadNotifications: &enableFileUploadNotifications,
			CloudToDevice:                 cloudToDeviceProperties,
		},
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, ok := d.GetOk("network_rule_set"); ok {
		props.Properties.NetworkRuleSets = expandNetworkRuleSetProperties(d)
	}

	//nolint:staticcheck
	if v, ok := d.GetOkExists("public_network_access_enabled"); ok {
		enabled := devices.PublicNetworkAccessDisabled
		if v.(bool) {
			enabled = devices.PublicNetworkAccessEnabled
		}
		props.Properties.PublicNetworkAccess = enabled
	}

	retention, retentionOk := d.GetOk("event_hub_retention_in_days")
	partition, partitionOk := d.GetOk("event_hub_partition_count")
	if partitionOk || retentionOk {
		eh := devices.EventHubProperties{}
		if retentionOk {
			eh.RetentionTimeInDays = pointer.To(int64(retention.(int)))
		}
		if partitionOk {
			eh.PartitionCount = pointer.To(int32(partition.(int)))
		}

		props.Properties.EventHubEndpoints = map[string]*devices.EventHubProperties{
			"events": &eh,
		}
	}

	props.Properties.DisableLocalAuth = pointer.To(!d.Get("local_authentication_enabled").(bool))

	if v, ok := d.GetOk("min_tls_version"); ok {
		props.Properties.MinTLSVersion = pointer.To(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, props, "")
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %q: %+v", id, err)
	}

	return resourceIotHubRead(d, meta)
}
