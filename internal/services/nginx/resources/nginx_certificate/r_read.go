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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxcertificate.ParseCertificateID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Nginx.NginxCertificate
			result, err := client.CertificatesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}
			var output CertificateModel

			output.Name = pointer.From(result.Model.Name)
			output.NginxDeploymentId = nginxdeployment.NewNginxDeploymentID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName).ID()
			prop := result.Model.Properties
			output.KeyVirtualPath = pointer.From(prop.KeyVirtualPath)
			output.KeyVaultSecretId = pointer.From(prop.KeyVaultSecretId)
			output.CertificateVirtualPath = pointer.From(prop.CertificateVirtualPath)
			return meta.Encode(&output)
		},
	}
}
