// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.CertificateClient_v2024_05_01

			var model ApiManagementWorkspaceCertificateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := certificate.NewWorkspaceCertificateID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.WorkspaceCertificateGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := certificate.CertificateCreateOrUpdateParameters{
				Properties: &certificate.CertificateCreateOrUpdateProperties{},
			}

			if model.KeyVaultSecretId != "" {
				parsedSecretId, err := keyvault.ParseNestedItemID(model.KeyVaultSecretId, keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret)
				if err != nil {
					return err
				}

				parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{
					SecretIdentifier: pointer.To(parsedSecretId.ID()),
				}

				if model.UserAssignedIdentityClientId != "" {
					parameters.Properties.KeyVault.IdentityClientId = pointer.To(model.UserAssignedIdentityClientId)
				}
			}

			if model.CertificateDataBase64 != "" {
				parameters.Properties.Data = pointer.To(model.CertificateDataBase64)
				parameters.Properties.Password = pointer.To(model.Password)
			}

			if _, err := client.WorkspaceCertificateCreateOrUpdate(ctx, id, parameters, certificate.DefaultWorkspaceCertificateCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
