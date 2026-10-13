// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.CertificateClient_v2024_05_01

			var model ApiManagementWorkspaceCertificateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := certificate.ParseWorkspaceCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.WorkspaceCertificateGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := resp.Model
			parameters := certificate.CertificateCreateOrUpdateParameters{
				Properties: &certificate.CertificateCreateOrUpdateProperties{},
			}

			if kv := payload.Properties.KeyVault; kv != nil {
				parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{
					SecretIdentifier: kv.SecretIdentifier,
					IdentityClientId: kv.IdentityClientId,
				}
			}

			if metadata.ResourceData.HasChange("user_assigned_identity_client_id") {
				if model.UserAssignedIdentityClientId != "" {
					if parameters.Properties.KeyVault == nil {
						parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{}
					}
					parameters.Properties.KeyVault.IdentityClientId = pointer.To(model.UserAssignedIdentityClientId)
				} else if parameters.Properties.KeyVault != nil {
					parameters.Properties.KeyVault.IdentityClientId = nil
				}
			}

			if metadata.ResourceData.HasChange("key_vault_secret_id") {
				if model.KeyVaultSecretId != "" {
					parsedSecretId, err := keyvault.ParseNestedItemID(model.KeyVaultSecretId, keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret)
					if err != nil {
						return err
					}

					if parameters.Properties.KeyVault == nil {
						parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{}
					}
					parameters.Properties.KeyVault.SecretIdentifier = pointer.To(parsedSecretId.ID())
				} else {
					parameters.Properties.KeyVault = nil
				}
			}

			if metadata.ResourceData.HasChange("certificate_data_base64") {
				if model.CertificateDataBase64 != "" {
					parameters.Properties.Data = pointer.To(model.CertificateDataBase64)
				} else {
					parameters.Properties.Data = nil
				}
			}

			if metadata.ResourceData.HasChange("password") {
				if model.Password != "" {
					parameters.Properties.Password = pointer.To(model.Password)
				} else {
					parameters.Properties.Password = nil
				}
			}

			if _, err := client.WorkspaceCertificateCreateOrUpdate(ctx, *id, parameters, certificate.DefaultWorkspaceCertificateCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
