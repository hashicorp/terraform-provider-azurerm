// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_endpoint_servicebus_topic

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_endpoint_servicebus_topic/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// TODO: outside of this pr make this private

var IothubResourceName = "azurerm_iothub"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubEndpointServiceBusTopicCreateUpdate,
		Read:   resourceIotHubEndpointServiceBusTopicRead,
		Update: resourceIotHubEndpointServiceBusTopicCreateUpdate,
		Delete: resourceIotHubEndpointServiceBusTopicDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.IoTHubEndpointServiceBusTopicV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EndpointServiceBusTopicID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceIothubEndpointServicebusTopicSchema(),
	}
}
