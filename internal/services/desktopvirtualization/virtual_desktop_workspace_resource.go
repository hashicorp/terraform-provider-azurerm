// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/scalingplan"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource                   = DesktopVirtualizationWorkspaceResource{}
	_ sdk.ResourceWithUpdate         = DesktopVirtualizationWorkspaceResource{}
	_ sdk.ResourceWithStateMigration = DesktopVirtualizationWorkspaceResource{}
)

type DesktopVirtualizationWorkspaceResource struct{}

func (DesktopVirtualizationWorkspaceResource) ModelObject() interface{} {
	return &DesktopVirtualizationWorkspaceModel{}
}

func (DesktopVirtualizationWorkspaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scalingplan.ValidateHostPoolID
}

func (DesktopVirtualizationWorkspaceResource) ResourceType() string {
	return "azurerm_virtual_desktop_workspace"
}

func (DesktopVirtualizationWorkspaceResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.WorkspaceV0ToV1{},
		},
	}
}

type DesktopVirtualizationWorkspaceModel struct {
	Name                string            `tfschema:"name"`
	ResourceGroup       string            `tfschema:"resource_group_name"`
	Location            string            `tfschema:"location"`
	FriendlyName        string            `tfschema:"friendly_name"`
	Description         string            `tfschema:"description"`
	PublicNetworkAccess bool              `tfschema:"public_network_access_enabled"`
	Tags                map[string]string `tfschema:"tags"`
}

func (r DesktopVirtualizationWorkspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WorkspaceName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r DesktopVirtualizationWorkspaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DesktopVirtualizationWorkspaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.WorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DesktopVirtualizationWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			log.Printf("[INFO] preparing arguments for Virtual Desktop Workspace create")

			id := workspace.NewWorkspaceID(subscriptionId, model.ResourceGroup, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := workspace.Workspace{
				Location: azure.NormalizeLocation(model.Location),
				Tags:     pointer.To(model.Tags),
				Properties: &workspace.WorkspaceProperties{
					Description:  pointer.To(model.Description),
					FriendlyName: pointer.To(model.FriendlyName),
				},
			}

			publicNetworkAccess := workspace.PublicNetworkAccessEnabled

			if !model.PublicNetworkAccess {
				publicNetworkAccess = workspace.PublicNetworkAccessDisabled
			}

			payload.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r DesktopVirtualizationWorkspaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.WorkspacesClient

			state := DesktopVirtualizationWorkspaceModel{}

			id, err := workspace.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state.Name = id.WorkspaceName
			state.ResourceGroup = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.FriendlyName = pointer.From(props.FriendlyName)
					publicNetworkAccess := true
					if v := props.PublicNetworkAccess; v != nil && *v != workspace.PublicNetworkAccessEnabled {
						publicNetworkAccess = false
					}
					state.PublicNetworkAccess = publicNetworkAccess
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DesktopVirtualizationWorkspaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.WorkspacesClient

			var model DesktopVirtualizationWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			log.Printf("[INFO] preparing arguments for Virtual Desktop Workspace update")

			id, err := workspace.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("friendly_name") {
				payload.Properties.FriendlyName = pointer.To(model.FriendlyName)
			}

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			publicNetworkAccess := workspace.PublicNetworkAccessEnabled

			if !model.PublicNetworkAccess {
				publicNetworkAccess = workspace.PublicNetworkAccessDisabled
			}

			payload.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)

			if _, err := client.CreateOrUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DesktopVirtualizationWorkspaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.WorkspacesClient

			id, err := workspace.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.WorkspaceName, r.ResourceType())
			defer locks.UnlockByName(id.WorkspaceName, r.ResourceType())

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
