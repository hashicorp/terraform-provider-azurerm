// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_namespace/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var (
	serviceBusNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"
	serviceBusNamespaceResourceName             = "azurerm_servicebus_namespace"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a ServiceBus Namespace.",
		Create:      resourceServiceBusNamespaceCreate,
		Read:        resourceServiceBusNamespaceRead,
		Update:      resourceServiceBusNamespaceUpdate,
		Delete:      resourceServiceBusNamespaceDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&namespaces.NamespaceId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&namespaces.NamespaceId{}),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NamespaceV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: servicebusNamespaceSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v any) error {
				oldCustomerManagedKey, newCustomerManagedKey := diff.GetChange("customer_managed_key")
				if len(oldCustomerManagedKey.([]any)) != 0 && len(newCustomerManagedKey.([]any)) == 0 {
					diff.ForceNew("customer_managed_key")
				}

				oldSku, newSku := diff.GetChange("sku")
				if diff.HasChange("sku") {
					if strings.EqualFold(newSku.(string), string(namespaces.SkuNamePremium)) || strings.EqualFold(oldSku.(string), string(namespaces.SkuNamePremium)) {
						log.Printf("[DEBUG] cannot migrate a namespace from or to Premium SKU")
						diff.ForceNew("sku")
					}
				}
				return nil
			}),
			pluginsdk.CustomizeDiffShim(servicebusTLSVersionDiff),
		),
	}
}
