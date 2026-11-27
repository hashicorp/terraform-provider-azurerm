// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_network_rules

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -parent-id "storage_account_id"

var storageAccountResourceName = "azurerm_storage_account"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountNetworkRulesCreate,
		Read:   resourceStorageAccountNetworkRulesRead,
		Update: resourceStorageAccountNetworkRulesUpdate,
		Delete: resourceStorageAccountNetworkRulesDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.StorageAccountId{}, pluginsdk.ResourceTypeForIdentityVirtual),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.StorageAccountId{}, pluginsdk.ResourceTypeForIdentityVirtual),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: storageAccountNetworkRulesSchema(),
	}
}
