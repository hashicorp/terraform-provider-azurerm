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

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxApiKey
			id, err := nginxapikey.ParseApiKeyID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ApiKeysGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: got nil model", id)
			}

			var state APIKeyModel
			if err := meta.Decode(&state); err != nil {
				return err
			}

			var output APIKeyModel
			output.Name = id.ApiKeyName
			output.NginxDeploymentId = nginxdeployment.NewNginxDeploymentID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName).ID()
			// secret_text field is not returned by the API so decode from state
			output.SecretText = state.SecretText
			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					output.EndDateTime = pointer.From(props.EndDateTime)
					output.Hint = pointer.From(props.Hint)
				}
			}
			return meta.Encode(&output)
		},
	}
}
