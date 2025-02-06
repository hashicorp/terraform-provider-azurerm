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
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DesktopVirtualizationWorkspaceDataSource struct{}

type DesktopVirtualizationWorkspaceModel struct {
	Name                string            `tfschema:"name"`
	ResourceGroup       string            `tfschema:"resource_group_name"`
	Location            string            `tfschema:"location"`
	FriendlyName        string            `tfschema:"friendly_name"`
	Description         string            `tfschema:"description"`
	PublicNetworkAccess bool              `tfschema:"public_network_access_enabled"`
	Tags                map[string]string `tfschema:"tags"`
}

var _ sdk.DataSource = DesktopVirtualizationWorkspaceDataSource{}

func (d DesktopVirtualizationWorkspaceDataSource) ModelObject() interface{} {
	return &DesktopVirtualizationWorkspaceModel{}
}

func (d DesktopVirtualizationWorkspaceDataSource) ResourceType() string {
	return "azurerm_virtual_desktop_workspace"
}

func (d DesktopVirtualizationWorkspaceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WorkspaceName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d DesktopVirtualizationWorkspaceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
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

func (d DesktopVirtualizationWorkspaceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.WorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DesktopVirtualizationWorkspaceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id := workspace.NewWorkspaceID(subscriptionId, state.ResourceGroup, state.Name)

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

			state.Location = location.NormalizeNilable(model.Location)
			state.Tags = pointer.From(model.Tags)

			if properties := model.Properties; properties != nil {
				if properties.FriendlyName != nil {
					state.FriendlyName = pointer.From(properties.FriendlyName)
				}

				if properties.Description != nil {
					state.Description = pointer.From(properties.Description)
				}

				publicNetworkAccess := true
				if properties.PublicNetworkAccess != nil && pointer.From(properties.PublicNetworkAccess) != workspace.PublicNetworkAccessEnabled {
					publicNetworkAccess = false
				}
				state.PublicNetworkAccess = publicNetworkAccess
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
