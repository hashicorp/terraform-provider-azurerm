// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "namespace_name:namespace_id,resource_group_name:namespace_id,subscription_id:namespace_id"

const serviceBusQueueResourceName = "azurerm_servicebus_queue"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a ServiceBus Queue.",
		Create:      resourceServiceBusQueueCreateUpdate,
		Read:        resourceServiceBusQueueRead,
		Update:      resourceServiceBusQueueCreateUpdate,
		Delete:      resourceServiceBusQueueDelete,

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&queues.QueueId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&queues.QueueId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceServicebusQueueSchema(),
	}
}
