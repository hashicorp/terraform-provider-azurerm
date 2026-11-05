// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DiskEncryptionSetName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
		},

		"auto_key_rotation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"encryption_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(diskencryptionsets.DiskEncryptionSetTypeEncryptionAtRestWithCustomerKey),
			ValidateFunc: validation.StringInSlice(diskencryptionsets.PossibleValuesForDiskEncryptionSetType(), false),
		},

		"federated_client_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"tags": commonschema.Tags(),

		"key_vault_key_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
