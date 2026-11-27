// Copyright IBM Corp. 2014, 2026
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type APIKeyDataSourceModel struct {
	Name              string `tfschema:"name"`
	NginxDeploymentId string `tfschema:"nginx_deployment_id"`
	Hint              string `tfschema:"hint"`
	EndDateTime       string `tfschema:"end_date_time"`
}

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

func (m DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the NGINX Dataplane API Key.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"nginx_deployment_id": {
			Description:  "The ID of the NGINX Deployment that the API key is associated with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: nginxdeployment.ValidateNginxDeploymentID,
		},
	}
}

func (m DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"end_date_time": {
			Description: "The RFC3339 formatted time after which this Dataplane API Key is no longer valid.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"hint": {
			Description: "The first three characters of the secret text to help identify it in use.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (m DataSource) ModelObject() any {
	return &APIKeyDataSourceModel{}
}

func (m DataSource) ResourceType() string {
	return "azurerm_nginx_api_key"
}

func (m DataSource) Description() string {
	return "Use this data source to access information about an existing NGINX Dataplane API Key."
}

func (m DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxApiKey
			var state APIKeyDataSourceModel
			if err := meta.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			deploymentId, err := nginxdeployment.ParseNginxDeploymentID(state.NginxDeploymentId)
			if err != nil {
				return err
			}
			id := nginxapikey.NewApiKeyID(
				deploymentId.SubscriptionId,
				deploymentId.ResourceGroupName,
				deploymentId.NginxDeploymentName,
				state.Name,
			)

			resp, err := client.ApiKeysGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil && model.Properties != nil {
				props := model.Properties
				state.EndDateTime = pointer.From(props.EndDateTime)
				state.Hint = pointer.From(props.Hint)
			}

			meta.SetID(id)
			return meta.Encode(&state)
		},
	}
}
