// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthbot

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2025-05-25/healthbots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthbotServiceCreate,
		Read:   resourceHealthbotServiceRead,
		Update: resourceHealthbotServiceUpdate,
		Delete: resourceHealthbotServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := healthbots.ParseHealthBotID(id)
			return err
		}),

		Schema: healthbotSchema(),

		CustomizeDiff: func(_ context.Context, diff *pluginsdk.ResourceDiff, _ any) error {
			if diff.HasChange("sku_name") {
				// Downgrading to `F0` isn't possible, and the API returns an uninformative error
				if _, newSKU := diff.GetChange("sku_name"); newSKU == string(healthbots.SkuNameFZero) {
					if err := diff.ForceNew("sku_name"); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
