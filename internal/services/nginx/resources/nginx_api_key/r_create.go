// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_api_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxapikey"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Nginx.NginxApiKey

			var model APIKeyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			deployID, err := nginxdeployment.ParseNginxDeploymentID(model.NginxDeploymentId)
			if err != nil {
				return err
			}

			subscriptionID := metadata.Client.Account.SubscriptionId
			id := nginxapikey.NewApiKeyID(
				subscriptionID,
				deployID.ResourceGroupName,
				deployID.NginxDeploymentName,
				model.Name,
			)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ApiKeysGet(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %+v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			req := nginxapikey.NginxDeploymentApiKeyRequest{
				Properties: &nginxapikey.NginxDeploymentApiKeyRequestProperties{
					SecretText:  pointer.To(model.SecretText),
					EndDateTime: pointer.To(model.EndDateTime),
				},
			}

			if _, err = client.ApiKeysCreateOrUpdate(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
