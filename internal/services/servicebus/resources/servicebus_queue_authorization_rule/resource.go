// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_queue_authorization_rule

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusQueueAuthorizationRuleCreateUpdate,
		Read:   resourceServiceBusQueueAuthorizationRuleRead,
		Update: resourceServiceBusQueueAuthorizationRuleCreateUpdate,
		Delete: resourceServiceBusQueueAuthorizationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := queues.ParseQueueAuthorizationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceServiceBusqueueAuthorizationRuleSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(helpers.AuthorizationRuleCustomizeDiff),
	}
}
