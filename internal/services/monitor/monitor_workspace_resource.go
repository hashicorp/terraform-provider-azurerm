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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WorkspaceResourceModel struct {
	Name                       string            `tfschema:"name"`
	ResourceGroupName          string            `tfschema:"resource_group_name"`
	PublicNetworkAccessEnabled bool              `tfschema:"public_network_access_enabled"`
	Location                   string            `tfschema:"location"`
	Tags                       map[string]string `tfschema:"tags"`
}

type WorkspaceResource struct{}

var _ sdk.ResourceWithUpdate = WorkspaceResource{}

func (r WorkspaceResource) ResourceType() string {
	return "azurerm_monitor_workspace"
}

func (r WorkspaceResource) ModelObject() interface{} {
	return &WorkspaceResourceModel{}
}

func (r WorkspaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return azuremonitorworkspaces.ValidateAccountID
}

func (r WorkspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r WorkspaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkspaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Monitor.WorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := azuremonitorworkspaces.NewAccountID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			publicNetworkAccess := azuremonitorworkspaces.PublicNetworkAccessEnabled
			if !model.PublicNetworkAccessEnabled {
				publicNetworkAccess = azuremonitorworkspaces.PublicNetworkAccessDisabled
			}

			properties := azuremonitorworkspaces.AzureMonitorWorkspaceResource{
				Location: location.Normalize(model.Location),
				Properties: &azuremonitorworkspaces.AzureMonitorWorkspace{
					PublicNetworkAccess: pointer.To(publicNetworkAccess),
				},
				Tags: &model.Tags,
			}

			if _, err := client.Create(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkspaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.WorkspacesClient

			id, err := azuremonitorworkspaces.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkspaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				resp.Model.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if properties.Properties == nil {
					properties.Properties = &azuremonitorworkspaces.AzureMonitorWorkspace{}
				}

				publicNetworkAccess := azuremonitorworkspaces.PublicNetworkAccessEnabled
				if !model.PublicNetworkAccessEnabled {
					publicNetworkAccess = azuremonitorworkspaces.PublicNetworkAccessDisabled
				}

				resp.Model.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
			}

			if _, err := client.Create(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r WorkspaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.WorkspacesClient

			id, err := azuremonitorworkspaces.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := WorkspaceResourceModel{
				Name:              id.AccountName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				state.Location = location.Normalize(model.Location)

				if properties := model.Properties; properties != nil {
					publicNetworkAccess := true
					if properties.PublicNetworkAccess != nil {
						publicNetworkAccess = azuremonitorworkspaces.PublicNetworkAccessEnabled == *properties.PublicNetworkAccess
					}
					state.PublicNetworkAccessEnabled = publicNetworkAccess
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkspaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.WorkspacesClient

			id, err := azuremonitorworkspaces.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
