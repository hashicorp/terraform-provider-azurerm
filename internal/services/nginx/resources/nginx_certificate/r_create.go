// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxcertificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	nginxClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/client"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := nginxClient.From(metadata.Client).NginxCertificate

			var model CertificateModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			deployID, _ := nginxdeployment.ParseNginxDeploymentID(model.NginxDeploymentId)

			subscriptionID := account.From(metadata.Client).SubscriptionId
			id := nginxcertificate.NewCertificateID(subscriptionID, deployID.ResourceGroupName, deployID.NginxDeploymentName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.CertificatesGet(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			req := nginxcertificate.NginxCertificate{
				Properties: &nginxcertificate.NginxCertificateProperties{
					CertificateVirtualPath: pointer.To(model.CertificateVirtualPath),
					KeyVaultSecretId:       pointer.To(model.KeyVaultSecretId),
					KeyVirtualPath:         pointer.To(model.KeyVirtualPath),
				},
			}

			if err := client.CertificatesCreateOrUpdateCallbackThenPoll(ctx, id, req, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
