// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_endpoint_eventhub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubEndpointEventHubCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).IoTHub.ResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	endpointRG := d.Get("resource_group_name").(string)

	iotHubId, err := parse.IotHubID(d.Get("iothub_id").(string))
	if err != nil {
		return err
	}
	iotHubName := iotHubId.Name
	iotHubRG := iotHubId.ResourceGroup

	id := parse.NewEndpointEventhubID(subscriptionId, iotHubRG, iotHubName, d.Get("name").(string))

	locks.ByName(iotHubName, IothubResourceName)
	defer locks.UnlockByName(iotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, iotHubRG, iotHubName)
	if err != nil {
		if response.WasNotFound(iothub.Response.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iotHubName, iotHubRG)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", iotHubName, iotHubRG, err)
	}

	authenticationType := devices.AuthenticationType(d.Get("authentication_type").(string))

	eventhubEndpoint := devices.RoutingEventHubProperties{
		AuthenticationType: authenticationType,
		Name:               pointer.To(id.EndpointName),
		ResourceGroup:      pointer.To(endpointRG),
	}

	// To align with the previous TF behaviour, `subscription_id` needs to be set with the provider's subscription Id when it isn't specified in the tf config, otherwise TF behaviour is different than before and it may block the existing users
	// From the business perspective, the raw config handling is only meant for the case that the user has an CosmosDB Account whose Endpoint's subscription is not the provider's one. Then the user wants to reset it to the provider's one by unset the subscription_id
	// From the TF code perspective, given `Computed: true` is enabled, TF would always get the value from the last apply when this property isn't set in the tf config. So `d.GetRawConfig()` is required to determine if it's set in the tf config
	if v := d.GetRawConfig().AsValueMap()["subscription_id"]; v.IsNull() {
		eventhubEndpoint.SubscriptionID = pointer.To(meta.(*clients.Client).Account.SubscriptionId)
	} else {
		eventhubEndpoint.SubscriptionID = pointer.To(d.Get("subscription_id").(string))
	}

	if authenticationType == devices.AuthenticationTypeKeyBased {
		if v, ok := d.GetOk("connection_string"); ok {
			eventhubEndpoint.ConnectionString = pointer.To(v.(string))
		} else {
			return fmt.Errorf("`connection_string` must be specified when `authentication_type` is `keyBased`")
		}
	} else {
		if v, ok := d.GetOk("endpoint_uri"); ok {
			eventhubEndpoint.EndpointURI = pointer.To(v.(string))
			eventhubEndpoint.EntityPath = pointer.To(d.Get("entity_path").(string))
		} else {
			return fmt.Errorf("`endpoint_uri` and `entity_path` must be specified when `authentication_type` is `identityBased`")
		}

		if v, ok := d.GetOk("identity_id"); ok {
			eventhubEndpoint.Identity = &devices.ManagedIdentity{
				UserAssignedIdentity: pointer.To(v.(string)),
			}
		}
	}

	routing := iothub.Properties.Routing
	if routing == nil {
		routing = &devices.RoutingProperties{}
	}

	if routing.Endpoints == nil {
		routing.Endpoints = &devices.RoutingEndpoints{}
	}

	if routing.Endpoints.EventHubs == nil {
		eventHubs := make([]devices.RoutingEventHubProperties, 0)
		routing.Endpoints.EventHubs = &eventHubs
	}

	endpoints := make([]devices.RoutingEventHubProperties, 0)

	alreadyExists := false
	for _, existingEndpoint := range *routing.Endpoints.EventHubs {
		if existingEndpointName := existingEndpoint.Name; existingEndpointName != nil {
			if strings.EqualFold(*existingEndpointName, id.EndpointName) {
				if d.IsNewResource() {
					if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError("azurerm_iothub_endpoint_eventhub", id.ID())
					}
				}
				endpoints = append(endpoints, eventhubEndpoint)
				alreadyExists = true
			} else {
				endpoints = append(endpoints, existingEndpoint)
			}
		}
	}

	if d.IsNewResource() {
		endpoints = append(endpoints, eventhubEndpoint)
	} else if !alreadyExists {
		return fmt.Errorf("unable to find %s", id)
	}
	routing.Endpoints.EventHubs = &endpoints

	future, err := client.CreateOrUpdate(ctx, iotHubRG, iotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	return resourceIotHubEndpointEventHubRead(d, meta)
}
