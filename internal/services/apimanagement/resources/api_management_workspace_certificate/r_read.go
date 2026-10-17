// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.CertificateClient_v2024_05_01

			id, err := certificate.ParseWorkspaceCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceCertificateGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := ApiManagementWorkspaceCertificateModel{
				Name:                     id.CertificateId,
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if props := respModel.Properties; props != nil {
					model.Expiration = props.ExpirationDate
					model.Subject = props.Subject
					model.Thumbprint = props.Thumbprint
					// The API omitted `certificate_data_base64` from the response as it is considered sensitive.
					model.CertificateDataBase64 = metadata.ResourceData.Get("certificate_data_base64").(string)
					// The API omitted `password` from the response as it is considered sensitive.
					model.Password = metadata.ResourceData.Get("password").(string)

					if kv := props.KeyVault; kv != nil {
						model.KeyVaultSecretId = pointer.From(kv.SecretIdentifier)
						model.UserAssignedIdentityClientId = pointer.From(kv.IdentityClientId)
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}
