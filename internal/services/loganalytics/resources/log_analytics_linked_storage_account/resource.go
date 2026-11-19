// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_linked_storage_account

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedstorageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/resources/log_analytics_linked_storage_account/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsLinkedStorageAccountCreateUpdate,
		Read:   resourceLogAnalyticsLinkedStorageAccountRead,
		Update: resourceLogAnalyticsLinkedStorageAccountCreateUpdate,
		Delete: resourceLogAnalyticsLinkedStorageAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := linkedstorageaccounts.ParseDataSourceTypeID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.LinkedStorageAccountV0ToV1{},
		}),

		Schema: logAnalyticsLinkedStorageAccountSchema(),
	}
}
