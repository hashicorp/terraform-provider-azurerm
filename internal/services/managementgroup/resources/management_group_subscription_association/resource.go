// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_subscription_association

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/management/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/resources/management_group_subscription_association/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var managementGroupCacheControl = "no-cache"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagementGroupSubscriptionAssociationCreate,
		Read:   resourceManagementGroupSubscriptionAssociationRead,
		Delete: resourceManagementGroupSubscriptionAssociationDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ManagementGroupSubscriptionAssociationV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managementgroups.ParseSubscriptionID(id)
			return err
		}),

		Schema: managementGroupSubscriptionAssociationSchema(),
	}
}
