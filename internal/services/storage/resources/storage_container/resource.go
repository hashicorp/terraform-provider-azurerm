// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_container/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:storage_account_id,resource_group_name:storage_account_id,storage_account_name:storage_account_id"

const StorageContainerResourceName = "azurerm_storage_container"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageContainerCreate,
		Read:   resourceStorageContainerRead,
		Delete: resourceStorageContainerDelete,
		Update: resourceStorageContainerUpdate,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.StorageContainerId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.StorageContainerId{}),
		},

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ContainerV0ToV1{},
			1: migration.StorageContainerV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: storageContainerSchema(),
	}
}
