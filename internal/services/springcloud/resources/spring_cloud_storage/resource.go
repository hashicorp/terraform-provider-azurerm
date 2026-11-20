// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_storage

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: "Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_storage` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.",

		Create: resourceSpringCloudStorageCreateUpdate,
		Read:   resourceSpringCloudStorageRead,
		Update: resourceSpringCloudStorageCreateUpdate,
		Delete: resourceSpringCloudStorageDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudCloudStorageV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
			// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
			_, err := appplatform.ParseStorageIDInsensitively(id)
			return err
		}),

		Schema: springCloudStorageSchema(),
	}
}
