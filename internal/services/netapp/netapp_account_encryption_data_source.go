// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetAppAccountEncryptionDataSource struct{}

var _ sdk.DataSource = NetAppAccountEncryptionDataSource{}

func (r NetAppAccountEncryptionDataSource) ModelObject() interface{} {
	return &netAppModels.NetAppAccountEncryptionDataSourceModel{}
}

func (r NetAppAccountEncryptionDataSource) ResourceType() string {
	return "azurerm_netapp_account_encryption"
}

func (r NetAppAccountEncryptionDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return netappaccounts.ValidateNetAppAccountID
}

func (r NetAppAccountEncryptionDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"netapp_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Description:  "The ID of the NetApp Account where encryption will be set.",
			ValidateFunc: netAppValidate.ValidateNetAppAccountID,
		},

		"user_assigned_identity_id": {
			Type:          pluginsdk.TypeString,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			Optional:      true,
			Description:   "The resource ID of the User Assigned Identity to use for encryption.",
			ConflictsWith: []string{"system_assigned_identity_principal_id"},
		},

		"system_assigned_identity_principal_id": {
			Type:          pluginsdk.TypeString,
			ValidateFunc:  validation.IsUUID,
			Optional:      true,
			Description:   "The Principal ID of the System Assigned Identity to use for encryption.",
			ConflictsWith: []string{"user_assigned_identity_id"},
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
					},
				},
			},
		},
	}
}

func (r NetAppAccountEncryptionDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetAppAccountEncryptionDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.AccountClient

			var state netAppModels.NetAppAccountEncryptionDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := netappaccounts.ParseNetAppAccountID(state.NetAppAccountID)
			if err != nil {
				return fmt.Errorf("error parsing netapp account id %s: %+v", state.NetAppAccountID, err)
			}

			resp, err := client.AccountsGet(ctx, pointer.From(id))
			if err != nil {
				if resp.HttpResponse.StatusCode == http.StatusNotFound {
					return fmt.Errorf("not found %s: %v", id, err)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := resp.Model
			if model.Properties.Encryption == nil {
				return fmt.Errorf("encryption information does not exist for %s", id)
			}

			anfAccountIdentityFlattened, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
			if err != nil {
				return err
			}

			state.Encryption = flattenEncryption(model.Properties.Encryption)

			if len(anfAccountIdentityFlattened) > 0 {

				if anfAccountIdentityFlattened[0].Type == identity.TypeSystemAssigned {
					state.SystemAssignedIdentityPrincipalID = anfAccountIdentityFlattened[0].PrincipalId
				}

				if anfAccountIdentityFlattened[0].Type == identity.TypeUserAssigned {
					state.UserAssignedIdentityID = anfAccountIdentityFlattened[0].IdentityIds[0]
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
