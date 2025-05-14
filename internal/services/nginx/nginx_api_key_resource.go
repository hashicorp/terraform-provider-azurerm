// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxapikey"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type APIKeyModel struct {
	Name              string `tfschema:"name"`
	NginxDeploymentId string `tfschema:"nginx_deployment_id"`
	SecretText        string `tfschema:"secret_text"`
	Hint              string `tfschema:"hint"`
	EndDateTime       string `tfschema:"end_date_time"`
}

type APIKeyResource struct{}

var _ sdk.ResourceWithUpdate = (*APIKeyResource)(nil)

func (m APIKeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"nginx_deployment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: nginxdeployment.ValidateNginxDeploymentID,
		},

		"end_date_time": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validate.EndDateTime,
			DiffSuppressFunc: suppress.RFC3339Time,
		},

		"secret_text": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m APIKeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m APIKeyResource) ModelObject() interface{} {
	return &APIKeyModel{}
}

func (m APIKeyResource) ResourceType() string {
	return "azurerm_nginx_api_key"
}

func (m APIKeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxApiKey

			var model APIKeyModel
			if err := meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			deployID, err := nginxdeployment.ParseNginxDeploymentID(model.NginxDeploymentId)
			if err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := nginxapikey.NewApiKeyID(
				subscriptionID,
				deployID.ResourceGroupName,
				deployID.NginxDeploymentName,
				model.Name,
			)

			existing, err := client.ApiKeysGet(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := nginxapikey.NginxDeploymentApiKeyRequest{
				Properties: &nginxapikey.NginxDeploymentApiKeyRequestProperties{
					SecretText:  pointer.To(model.SecretText),
					EndDateTime: pointer.To(model.EndDateTime),
				},
			}

			_, err = client.ApiKeysCreateOrUpdate(ctx, id, req)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m APIKeyResource) Update() sdk.ResourceFunc {
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

func (m APIKeyResource) Read() sdk.ResourceFunc {
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
					output.EndDateTime = pointer.ToString(props.EndDateTime)
					output.Hint = pointer.ToString(props.Hint)
				}
			}
			return meta.Encode(&output)
		},
	}
}

func (m APIKeyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxApiKey
			id, err := nginxapikey.ParseApiKeyID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.ApiKeysDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (m APIKeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nginxapikey.ValidateApiKeyID
}
