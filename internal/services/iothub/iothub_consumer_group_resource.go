// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

func resourceIotHubConsumerGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubConsumerGroupCreate,
		Read:   resourceIotHubConsumerGroupRead,
		Delete: resourceIotHubConsumerGroupDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.IoTHubConsumerGroupV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ConsumerGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubConsumerGroupName,
			},

			"iothub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"eventhub_endpoint_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),
		},
	}
}

func resourceIotHubConsumerGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM IoTHub Consumer Group creation.")

	id := parse.NewConsumerGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("eventhub_endpoint_name").(string), d.Get("name").(string))

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	if d.IsNewResource() {
		existing, err := client.GetEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Consumer Group %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_iothub_consumer_group", id.ID())
		}
	}

	consumerGroupBody := devices.EventHubConsumerGroupBodyDescription{
		// The properties are currently undocumented. See also:
		// https://docs.microsoft.com/en-us/azure/templates/microsoft.devices/2021-03-31/iothubs/eventhubendpoints/consumergroups?tabs=json#eventhubconsumergroupname
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

	read, err := client.GetEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubConsumerGroupRead(d, meta)
}

func resourceIotHubConsumerGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConsumerGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request for %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("iothub_name", id.IotHubName)
	d.Set("eventhub_endpoint_name", id.EventHubEndpointName)
	d.Set("resource_group_name", id.ResourceGroup)

	return nil
}

func resourceIotHubConsumerGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConsumerGroupID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	resp, err := client.DeleteEventHubConsumerGroup(ctx, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}
