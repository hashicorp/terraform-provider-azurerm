// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_subscription

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_subscription/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a ServiceBus Subscription.",
		Create:      resourceServiceBusSubscriptionCreateUpdate,
		Read:        resourceServiceBusSubscriptionRead,
		Update:      resourceServiceBusSubscriptionCreateUpdate,
		Delete:      resourceServiceBusSubscriptionDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ServiceBusSubscriptionV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := subscriptions.ParseSubscriptions2ID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceServicebusSubscriptionSchema(),
	}
}
