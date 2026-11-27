// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var SubscriptionResourceName = "azurerm_subscription"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubscriptionCreate,
		Update: resourceSubscriptionUpdate,
		Read:   resourceSubscriptionRead,
		Delete: resourceSubscriptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.SubscriptionAliasID(id)
			return err
		}, importSubscriptionByAlias()),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: subscriptionSchema(),
	}
}
