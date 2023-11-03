// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetAppAccountEncryptionResource struct{}

var _ sdk.Resource = NetAppAccountEncryptionResource{}

func (r NetAppAccountEncryptionResource) ModelObject() interface{} {
	return &netAppModels.NetAppAccountEncryption{}
}

func (r NetAppAccountEncryptionResource) ResourceType() string {
	return "azurerm_netapp_account_encryption"
}

func (r NetAppAccountEncryptionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return netappaccounts.ValidateNetAppAccountID
}

func (r NetAppAccountEncryptionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"netapp_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Description:  "The ID of the NetApp Account where encryption will be set.",
			ValidateFunc: netAppValidate.ValidateNetAppAccountID,
		},

		"user_assigned_identity_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			Description:   "The resource ID of the User Assigned Identity to use for encryption.",
			ConflictsWith: []string{"system_assigned_identity_principal_id"},
		},

		"system_assigned_identity_principal_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.IsUUID,
			Description:   "The Principal ID of the System Assigned Identity to use for encryption.",
			ConflictsWith: []string{"user_assigned_identity_id"},
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
					},
				},
			},
		},
	}
}

func (r NetAppAccountEncryptionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetAppAccountEncryptionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.AccountClient
			keyVaultsClient := metadata.Client.KeyVault
			resourcesClient := metadata.Client.Resource

			var model netAppModels.NetAppAccountEncryption
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountID, err := netappaccounts.ParseNetAppAccountID(model.NetAppAccountID)
			if err != nil {
				return fmt.Errorf("error parsing account id %s: %+v", model.NetAppAccountID, err)
			}

			metadata.Logger.Infof("Import check for %s", accountID.ID())
			existing, err := client.AccountsGet(ctx, pointer.From(accountID))
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("not found %s: %s", accountID.ID(), err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				if existing.Model.Properties.Encryption != nil && existing.Model.Properties.Encryption.KeySource != nil && pointer.From(existing.Model.Properties.Encryption.KeySource) == netappaccounts.KeySourceMicrosoftPointKeyVault {

					return tf.ImportAsExistsError(r.ResourceType(), accountID.ID())
				}
			}

			update := netappaccounts.NetAppAccountPatch{
				Properties: &netappaccounts.AccountProperties{},
			}

			encryptionExpanded, err := expandEncryption(ctx, model.Encryption, keyVaultsClient, resourcesClient, pointer.To(model))
			if err != nil {
				return err
			}
			if encryptionExpanded != nil && pointer.From(encryptionExpanded.KeySource) != netappaccounts.KeySourceMicrosoftPointKeyVault {
				return fmt.Errorf("encryption settings should be passed only for keyvault key provider")
			}

			update.Properties.Encryption = encryptionExpanded

			if err := client.AccountsUpdateThenPoll(ctx, pointer.From(accountID), update); err != nil {
				return fmt.Errorf("updating %s: %+v", accountID, err)
			}

			metadata.SetID(accountID)

			return nil
		},
	}
}

func (r NetAppAccountEncryptionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.AccountClient
			keyVaultsClient := metadata.Client.KeyVault
			resourcesClient := metadata.Client.Resource

			id, err := netappaccounts.ParseNetAppAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppAccountEncryption
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("Updating %s", id)

			update := netappaccounts.NetAppAccountPatch{
				Properties: &netappaccounts.AccountProperties{},
			}

			if metadata.ResourceData.HasChange("user_assigned_identity_id") || metadata.ResourceData.HasChange("system_assigned_identity_principal_id") || metadata.ResourceData.HasChange("encryption") {
				encryptionExpanded, err := expandEncryption(ctx, state.Encryption, keyVaultsClient, resourcesClient, pointer.To(state))
				if err != nil {
					return err
				}

				update.Properties.Encryption = encryptionExpanded

				if err := client.AccountsUpdateThenPoll(ctx, pointer.From(id), update); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}

				metadata.SetID(id)
			}

			return nil
		},
	}
}

func (r NetAppAccountEncryptionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.AccountClient

			id, err := netappaccounts.ParseNetAppAccountID((metadata.ResourceData.Id()))
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppAccountEncryption
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.AccountsGet(ctx, pointer.From(id))
			if err != nil {
				if existing.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if existing.Model.Properties.Encryption == nil {
				return fmt.Errorf("encryption information does not exist for %s", id)
			}

			anfAccountIdentityFlattened, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(existing.Model.Identity)
			if err != nil {
				return err
			}

			model := netAppModels.NetAppAccountEncryption{
				NetAppAccountID: id.ID(),
				Encryption:      flattenEncryption(existing.Model.Properties.Encryption),
			}

			if len(anfAccountIdentityFlattened) > 0 {

				if anfAccountIdentityFlattened[0].Type == identity.TypeSystemAssigned {
					model.SystemAssignedIdentityPrincipalID = anfAccountIdentityFlattened[0].PrincipalId
				}

				if anfAccountIdentityFlattened[0].Type == identity.TypeUserAssigned {
					model.UserAssignedIdentityID = anfAccountIdentityFlattened[0].IdentityIds[0]
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}

func (r NetAppAccountEncryptionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.AccountClient

			id, err := netappaccounts.ParseNetAppAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppAccountEncryption
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("Updating %s", id)

			update := netappaccounts.NetAppAccountPatch{
				Properties: &netappaccounts.AccountProperties{},
			}

			update.Properties.Encryption = &netappaccounts.AccountEncryption{}

			if err := client.AccountsUpdateThenPoll(ctx, pointer.From(id), update); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandEncryption(ctx context.Context, input []netAppModels.NetAppAccountEncryptionModel, keyVaultsClient *keyVaultClient.Client, resourcesClient *resourcesClient.Client, model *netAppModels.NetAppAccountEncryption) (*netappaccounts.AccountEncryption, error) {
	defaultEnc := netappaccounts.AccountEncryption{
		KeySource: pointer.To(netappaccounts.KeySourceMicrosoftPointNetApp),
	}

	if len(input) == 0 {
		return &defaultEnc, nil
	}

	keyId, err := keyVaultParse.ParseOptionallyVersionedNestedKeyID(input[0].KeyVaultKeyID)
	if err != nil {
		return nil, fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
	}

	keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyId.KeyVaultBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieving the resource id the key vault at url %q: %s", keyId.KeyVaultBaseUrl, err)
	}

	parsedKeyVaultID, err := commonids.ParseKeyVaultID(pointer.From(keyVaultID))
	if err != nil {
		return nil, err
	}

	encryptionIdentity := &netappaccounts.EncryptionIdentity{}

	if model.SystemAssignedIdentityPrincipalID == "" && model.UserAssignedIdentityID != "" {
		encryptionIdentity = &netappaccounts.EncryptionIdentity{
			UserAssignedIdentity: pointer.To(model.UserAssignedIdentityID),
		}
	}

	encryptionProperty := netappaccounts.AccountEncryption{
		Identity:  encryptionIdentity,
		KeySource: pointer.To(netappaccounts.KeySourceMicrosoftPointKeyVault),
		KeyVaultProperties: &netappaccounts.KeyVaultProperties{
			KeyName:            keyId.Name,
			KeyVaultUri:        keyId.KeyVaultBaseUrl,
			KeyVaultResourceId: parsedKeyVaultID.ID(),
		},
	}

	return &encryptionProperty, nil
}

func flattenEncryption(encryptionProperties *netappaccounts.AccountEncryption) []netAppModels.NetAppAccountEncryptionModel {
	if encryptionProperties == nil || *encryptionProperties.KeySource == netappaccounts.KeySourceMicrosoftPointNetApp {
		return []netAppModels.NetAppAccountEncryptionModel{}
	}

	return []netAppModels.NetAppAccountEncryptionModel{
		{
			KeyVaultKeyID: encryptionProperties.KeyVaultProperties.KeyVaultUri + "keys/" + encryptionProperties.KeyVaultProperties.KeyName,
		},
	}
}
