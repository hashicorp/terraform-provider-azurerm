// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_subscription_pricing

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/resources/security_center_subscription_pricing/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages the Pricing Tier for Azure Security Center in the current subscription.",
		Create:      resourceSecurityCenterSubscriptionPricingCreate,
		Read:        resourceSecurityCenterSubscriptionPricingRead,
		Update:      resourceSecurityCenterSubscriptionPricingUpdate,
		Delete:      resourceSecurityCenterSubscriptionPricingDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := pricings.ParsePricingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SubscriptionPricingV0ToV1{},
		}),

		Schema: securityCenterSubscriptionPricingSchema(),
	}
}
