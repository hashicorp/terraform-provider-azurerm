// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_endpoint_servicebus_queue

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubEndpointServiceBusQueueRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointServiceBusQueueID(d.Id())
	if err != nil {
		return err
	}

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if response.WasNotFound(iothub.Response.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	d.Set("name", id.EndpointName)

	iotHubId := parse.NewIotHubID(id.SubscriptionId, id.ResourceGroup, id.IotHubName)
	d.Set("iothub_id", iotHubId.ID())

	if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
		d.SetId("")
		return nil
	}

	exist := false

	if endpoints := iothub.Properties.Routing.Endpoints.ServiceBusQueues; endpoints != nil {
		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, id.EndpointName) {
					exist = true
					d.Set("resource_group_name", endpoint.ResourceGroup)
					d.Set("subscription_id", pointer.From(endpoint.SubscriptionID))

					authenticationType := string(devices.AuthenticationTypeKeyBased)
					if string(endpoint.AuthenticationType) != "" {
						authenticationType = string(endpoint.AuthenticationType)
					}
					d.Set("authentication_type", authenticationType)

					d.Set("connection_string", pointer.From(endpoint.ConnectionString))

					d.Set("endpoint_uri", pointer.From(endpoint.EndpointURI))

					d.Set("entity_path", pointer.From(endpoint.EntityPath))

					identityId := ""
					if endpoint.Identity != nil && endpoint.Identity.UserAssignedIdentity != nil {
						identityId = *endpoint.Identity.UserAssignedIdentity
					}
					d.Set("identity_id", identityId)
				}
			}
		}
	}

	if !exist {
		d.SetId("")
	}

	return nil
}
