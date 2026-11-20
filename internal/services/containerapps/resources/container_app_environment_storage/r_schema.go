// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedEnvironmentStorageName,
			Description:  "The name for this Storage.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironmentsstorages.ValidateManagedEnvironmentID,
			Description:  "The ID of the Container App Environment to which this storage belongs.",
		},

		"account_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  storageValidate.StorageAccountName,
			ConflictsWith: []string{"nfs_server_url"},
			Description:   "The Azure Storage Account in which the Share to be used is located.",
		},

		"access_key": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			RequiredWith:  []string{"account_name"},
			ConflictsWith: []string{"access_key_vault_secret"},
			Description:   "The Storage Account Access Key.",
		},

		"access_key_vault_secret": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			RequiredWith:  []string{"account_name"},
			ConflictsWith: []string{"access_key"},
			Description:   "A Key Vault Secret reference containing the Storage Account Access Key.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.Any(
							commonids.ValidateUserAssignedIdentityID,
							validation.StringInSlice([]string{"System"}, false),
						),
						Description: "The identity used to access the Key Vault. Possible values are `System` or the ID of a User Assigned Identity.",
					},
					"key_vault_secret_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
						Description:  "The ID of the Key Vault Secret containing the Storage Account Access Key. Can be either a versioned or versionless ID.",
					},
				},
			},
		},

		"share_name": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the Azure Storage Share to use.",
		},

		"access_mode": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(managedenvironmentsstorages.PossibleValuesForAccessMode(), false),
			Description:  "The access mode to connect this storage to the Container App. Possible values include `ReadOnly` and `ReadWrite`.",
		},

		"nfs_server_url": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"account_name"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			rawConfig := metadata.ResourceDiff.GetRawConfig()
			if rawConfig.IsNull() || !rawConfig.IsKnown() {
				return nil
			}

			config := rawConfig.AsValueMap()
			if config["account_name"].IsNull() {
				return nil
			}

			accessKeySet := !config["access_key"].IsNull()
			keyVaultSecretSet := !config["access_key_vault_secret"].IsNull() && (!config["access_key_vault_secret"].IsKnown() || config["access_key_vault_secret"].LengthInt() > 0)
			if !accessKeySet && !keyVaultSecretSet {
				return fmt.Errorf("one of `access_key` or `access_key_vault_secret` must be specified when `account_name` is set")
			}

			return nil
		},
	}
}
