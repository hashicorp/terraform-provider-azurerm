// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WorkspaceDataSource struct{}

var _ sdk.DataSource = WorkspaceDataSource{}

func (d WorkspaceDataSource) ModelObject() interface{} {
	return &WorkspaceDataSource{}
}

func (d WorkspaceDataSource) ResourceType() string {
	return "azurerm_monitor_workspace"
}

func (d WorkspaceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d WorkspaceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"query_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d WorkspaceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.WorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state WorkspaceResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := azuremonitorworkspaces.NewAccountID(subscriptionId, state.ResourceGroupName, state.Name)
			metadata.Logger.Infof("retrieving %s", id)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var enablePublicNetWorkAccess bool
			var location, queryEndpoint string
			var tag map[string]string

			if model := resp.Model; model != nil {
				location = azure.NormalizeLocation(model.Location)
				tag = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					if props.PublicNetworkAccess != nil {
						enablePublicNetWorkAccess = azuremonitorworkspaces.PublicNetworkAccessEnabled == *props.PublicNetworkAccess
					}
					if props.Metrics != nil && props.Metrics.PrometheusQueryEndpoint != nil {
						queryEndpoint = *props.Metrics.PrometheusQueryEndpoint
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&WorkspaceResourceModel{
				Location:                   location,
				Name:                       id.AccountName,
				PublicNetworkAccessEnabled: enablePublicNetWorkAccess,
				QueryEndpoint:              queryEndpoint,
				ResourceGroupName:          id.ResourceGroupName,
				Tags:                       tag,
			})
		},
	}
}
