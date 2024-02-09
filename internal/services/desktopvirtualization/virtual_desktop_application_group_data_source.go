// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DesktopVirtualizationApplicationGroupDataSource struct{}

type DesktopVirtualizationApplicationGroupModel struct {
	ApplicationGroupName string            `tfschema:"name"`
	ResourceGroupName    string            `tfschema:"resource_group_name"`
	Location             string            `tfschema:"location"`
	ApplicationGroupType string            `tfschema:"type"`
	HostPoolId           string            `tfschema:"host_pool_id"`
	WorkspaceId          string            `tfschema:"workspace_id"`
	FriendlyName         string            `tfschema:"friendly_name"`
	Description          string            `tfschema:"description"`
	Tags                 map[string]string `tfschema:"tags"`
}

var _ sdk.DataSource = DesktopVirtualizationApplicationGroupDataSource{}

func (r DesktopVirtualizationApplicationGroupDataSource) ModelObject() interface{} {
	return &DesktopVirtualizationApplicationGroupModel{}
}

func (r DesktopVirtualizationApplicationGroupDataSource) ResourceType() string {
	return "azurerm_virtual_desktop_application_group"
}

func (r DesktopVirtualizationApplicationGroupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApplicationGroupName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r DesktopVirtualizationApplicationGroupDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"host_pool_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DesktopVirtualizationApplicationGroupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ApplicationGroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DesktopVirtualizationApplicationGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id := applicationgroup.NewApplicationGroupID(subscriptionId, state.ResourceGroupName, state.ApplicationGroupName)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state.ApplicationGroupName = id.ApplicationGroupName
			state.ResourceGroupName = id.ResourceGroupName
			state.Location = location.NormalizeNilable(model.Location)
			state.Tags = pointer.From(model.Tags)
			state.ApplicationGroupType = string(model.Properties.ApplicationGroupType)

			hostPoolId, err := hostpool.ParseHostPoolIDInsensitively(model.Properties.HostPoolArmPath)
			if err != nil {
				return fmt.Errorf("parsing Host Pool ID %q: %+v", model.Properties.HostPoolArmPath, err)
			}
			state.HostPoolId = hostPoolId.ID()

			if model.Properties.WorkspaceArmPath != nil {
				state.WorkspaceId = pointer.From(model.Properties.WorkspaceArmPath)
			}

			if model.Properties.FriendlyName != nil {
				state.FriendlyName = pointer.From(model.Properties.FriendlyName)
			}

			if model.Properties.Description != nil {
				state.Description = pointer.From(model.Properties.Description)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
