// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_access_policy

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var keyVaultResourceName = "azurerm_key_vault"

const keyVaultAccessPolicyResourceName = "azurerm_key_vault_access_policy"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a Key Vault Access Policy.",
		Create:      resourceKeyVaultAccessPolicyCreate,
		Read:        resourceKeyVaultAccessPolicyRead,
		Update:      resourceKeyVaultAccessPolicyUpdate,
		Delete:      resourceKeyVaultAccessPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AccessPolicyID(id)
			return err
		}),

		// NOTE: Resource Identity is implemented manually because `parse.AccessPolicyId` is a custom ID
		// type that does not implement `resourceids.ResourceId`, so the auto-generation helpers cannot be used.
		Identity: &schema.ResourceIdentity{
			SchemaFunc: func() map[string]*schema.Schema {
				return map[string]*schema.Schema{
					"key_vault_id": {
						Type:              schema.TypeString,
						RequiredForImport: true,
					},
					"object_id": {
						Type:              schema.TypeString,
						RequiredForImport: true,
					},
					"application_id": {
						Type:              schema.TypeString,
						RequiredForImport: true,
					},
				}
			},
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: keyVaultAccessPolicySchema(),
	}
}
