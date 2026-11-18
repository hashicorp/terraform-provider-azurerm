// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_consumer_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubConsumerGroupCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewConsumerGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("eventhub_endpoint_name").(string), d.Get("name").(string))

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Consumer Group %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.Response.Response) {
			return tf.ImportAsExistsError("azurerm_iothub_consumer_group", id.ID())
		}
	}

	consumerGroupBody := devices.EventHubConsumerGroupBodyDescription{
		// The properties are currently undocumented. See also:
		// https://docs.microsoft.com/azure/templates/microsoft.devices/2021-03-31/iothubs/eventhubendpoints/consumergroups?tabs=json#eventhubconsumergroupname
		//
		// There is an example where the name is repeated in the properties,
		// so that seems to be the "proper" way. See also:
		// https://github.com/Azure/azure-quickstart-templates/blob/master/quickstarts/microsoft.devices/iothub-with-consumergroup-create/azuredeploy.json#L74
		Properties: &devices.EventHubConsumerGroupName{
			Name: &id.Name,
		},
	}

	if _, err := client.CreateEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name, consumerGroupBody); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubConsumerGroupRead(d, meta)
}
