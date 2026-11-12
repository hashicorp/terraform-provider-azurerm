// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_api_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxapikey"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxApiKey
			id, err := nginxapikey.ParseApiKeyID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model APIKeyModel
			if err := meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.ApiKeysGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			// full update - fill in the existing fields from the API and then patch it
			upd := nginxapikey.NginxDeploymentApiKeyRequest{
				Name: existing.Model.Name,
				Properties: &nginxapikey.NginxDeploymentApiKeyRequestProperties{
					EndDateTime: existing.Model.Properties.EndDateTime,
					// secret_text field is not returned by the API so decode from state
					SecretText: pointer.To(model.SecretText),
				},
			}
			if meta.ResourceData.HasChange("end_date_time") {
				upd.Properties.EndDateTime = pointer.To(model.EndDateTime)
			}

			if _, err := client.ApiKeysCreateOrUpdate(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}
