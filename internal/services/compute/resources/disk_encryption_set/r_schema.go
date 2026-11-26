// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package disk_encryption_set

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func diskEncryptionSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Disk Encryption Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DiskEncryptionSetName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"key_vault_key_id": {
			Description:  "The URL to a Key Vault Key (either from a Key Vault Key, or the Key URL for the Key Vault Secret).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		},

		"auto_key_rotation_enabled": {
			Description: "Boolean flag to specify whether Azure Disk Encryption Set automatically rotates the encryption Key to latest version.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"encryption_type": {
			Description:  "The type of key used to encrypt the data of the disk.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(diskencryptionsets.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey),
			ValidateFunc: validation.StringInSlice(diskencryptionsets.PossibleValuesForDiskEncryptionSetType(), false),
		},

		"federated_client_id": {
			Description:  "Multi-tenant application client ID to access Key Vault in a different tenant.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"tags": commonschema.Tags(),

		"key_vault_key_url": {
			Description: "The URL for the Key Vault Key or Key Vault Secret that is currently being used by the service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
