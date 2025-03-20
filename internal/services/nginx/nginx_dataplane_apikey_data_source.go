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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataplaneAPIKeyDataSourceModel struct {
	Name              string `tfschema:"name"`
	NginxDeploymentId string `tfschema:"nginx_deployment_id"`
	Hint              string `tfschema:"hint"`
	EndDateTime       string `tfschema:"end_date_time"`
}

type DataplaneAPIKeyDataSource struct{}

var _ sdk.DataSource = DataplaneAPIKeyDataSource{}

func (m DataplaneAPIKeyDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"nginx_deployment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: nginxdeployment.ValidateNginxDeploymentID,
		},
	}
}

func (m DataplaneAPIKeyDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"end_date_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m DataplaneAPIKeyDataSource) ModelObject() interface{} {
	return &DataplaneAPIKeyDataSourceModel{}
}

func (m DataplaneAPIKeyDataSource) ResourceType() string {
	return "azurerm_nginx_dataplane_apikey"
}

func (m DataplaneAPIKeyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxApiKey
			var state DataplaneAPIKeyDataSourceModel
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
				state.EndDateTime = pointer.ToString(props.EndDateTime)
				state.Hint = pointer.ToString(props.Hint)
			}

			meta.SetID(id)
			return meta.Encode(&state)
		},
	}
}
