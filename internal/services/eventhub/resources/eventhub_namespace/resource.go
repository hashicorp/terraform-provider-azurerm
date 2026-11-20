// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/ehnamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	eventHubNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"
	eventHubNamespaceResourceName             = "azurerm_eventhub_namespace"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubNamespaceCreate,
		Read:   resourceEventHubNamespaceRead,
		Update: resourceEventHubNamespaceUpdate,
		Delete: resourceEventHubNamespaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := ehnamespaces.ParseNamespaceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: eventhubNamespaceSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
				oldSku, newSku := d.GetChange("sku")
				if d.HasChange("sku") {
					if strings.EqualFold(newSku.(string), string(ehnamespaces.SkuNamePremium)) || strings.EqualFold(oldSku.(string), string(ehnamespaces.SkuTierPremium)) {
						log.Printf("[DEBUG] cannot migrate a namespace from or to Premium SKU")
						d.ForceNew("sku")
					}
				}
				return nil
			}),
			pluginsdk.CustomizeDiffShim(eventhubTLSVersionDiff),
			pluginsdk.CustomizeDiffShim(eventhubGeoReplicationDiff),
		),
	}
}
