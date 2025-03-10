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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataplaneAPIKeyModel struct {
	Name              string `tfschema:"name"`
	NginxDeploymentId string `tfschema:"nginx_deployment_id"`
	SecretText        string `tfschema:"secret_text"`
	Hint              string `tfschema:"hint"`
	EndDateTime       string `tfschema:"end_date_time"`
}

type DataplaneAPIKeyResource struct{}

var _ sdk.ResourceWithUpdate = (*DataplaneAPIKeyResource)(nil)

func (m DataplaneAPIKeyResource) Arguments() map[string]*pluginsdk.Schema {
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

		"secret_text": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"end_date_time": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateFunc:     validation.IsRFC3339Time,
			DiffSuppressFunc: suppress.RFC3339Time,
		},
	}
}

func (m DataplaneAPIKeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m DataplaneAPIKeyResource) ModelObject() interface{} {
	return &DataplaneAPIKeyModel{}
}

func (m DataplaneAPIKeyResource) ResourceType() string {
	return "azurerm_nginx_dataplane_apikey"
}

func (m DataplaneAPIKeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxApiKey

			var model DataplaneAPIKeyModel
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
					return fmt.Errorf("retreiving %s: %+v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := nginxapikey.NginxDeploymentApiKeyRequest{
				Properties: &nginxapikey.NginxDeploymentApiKeyRequestProperties{
					SecretText: pointer.To(model.SecretText),
				},
			}
			if model.EndDateTime != "" {
				req.Properties.EndDateTime = pointer.To(model.EndDateTime)
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

func (m DataplaneAPIKeyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxApiKey
			id, err := nginxapikey.ParseApiKeyID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DataplaneAPIKeyModel
			if err := meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			param := nginxapikey.NginxDeploymentApiKeyRequest{
				Properties: &nginxapikey.NginxDeploymentApiKeyRequestProperties{
					SecretText: pointer.To(model.SecretText),
				},
			}
			if model.EndDateTime != "" {
				param.Properties.EndDateTime = pointer.To(model.EndDateTime)
			}
			if _, err := client.ApiKeysCreateOrUpdate(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (m DataplaneAPIKeyResource) Read() sdk.ResourceFunc {
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

			var state DataplaneAPIKeyModel
			if err := meta.Decode(&state); err != nil {
				return err
			}

			var output DataplaneAPIKeyModel
			output.Name = id.ApiKeyName
			output.NginxDeploymentId = nginxdeployment.NewNginxDeploymentID(id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName).ID()
			// secret_text field is not returned by the API so decode from state
			output.SecretText = state.SecretText
			if model := resp.Model; model != nil && model.Properties != nil {
				props := model.Properties
				output.EndDateTime = pointer.ToString(props.EndDateTime)
				output.Hint = pointer.ToString(props.Hint)
			}
			return meta.Encode(&output)
		},
	}
}

func (m DataplaneAPIKeyResource) Delete() sdk.ResourceFunc {
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

func (m DataplaneAPIKeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nginxapikey.ValidateApiKeyID
}
