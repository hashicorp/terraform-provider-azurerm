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
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/desktop"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource                   = DesktopVirtualizationApplicationGroupResource{}
	_ sdk.ResourceWithUpdate         = DesktopVirtualizationApplicationGroupResource{}
	_ sdk.ResourceWithStateMigration = DesktopVirtualizationApplicationGroupResource{}
)

type DesktopVirtualizationApplicationGroupResource struct{}

func (DesktopVirtualizationApplicationGroupResource) ModelObject() interface{} {
	return &DesktopVirtualizationApplicationGroupResourceModel{}
}

func (DesktopVirtualizationApplicationGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return applicationgroup.ValidateApplicationGroupID
}

func (DesktopVirtualizationApplicationGroupResource) ResourceType() string {
	return "azurerm_virtual_desktop_application_group"
}

type DesktopVirtualizationApplicationGroupResourceModel struct {
	Name                      string            `tfschema:"name"`
	Location                  string            `tfschema:"location"`
	ResourceGroupName         string            `tfschema:"resource_group_name"`
	Type                      string            `tfschema:"type"`
	HostPoolId                string            `tfschema:"host_pool_id"`
	FriendlyName              string            `tfschema:"friendly_name"`
	DefaultDesktopDisplayName string            `tfschema:"default_desktop_display_name"`
	Description               string            `tfschema:"description"`
	Tags                      map[string]string `tfschema:"tags"`
}

func (DesktopVirtualizationApplicationGroupResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ApplicationGroupV0ToV1{},
		},
	}
}

func (DesktopVirtualizationApplicationGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApplicationGroupName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(applicationgroup.PossibleValuesForApplicationGroupType(), false),
		},

		"host_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hostpool.ValidateHostPoolID,
		},

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},

		"default_desktop_display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"tags": commonschema.Tags(),
	}
}

func (DesktopVirtualizationApplicationGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DesktopVirtualizationApplicationGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ApplicationGroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DesktopVirtualizationApplicationGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			log.Printf("[INFO] preparing arguments for Virtual Desktop Application Group creation")

			locks.ByName(model.Name, r.ResourceType())
			defer locks.UnlockByName(model.Name, r.ResourceType())

			id := applicationgroup.NewApplicationGroupID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := applicationgroup.ApplicationGroup{
				Location: model.Location,
				Tags:     pointer.To(model.Tags),
				Properties: applicationgroup.ApplicationGroupProperties{
					ApplicationGroupType: applicationgroup.ApplicationGroupType(model.Type),
					FriendlyName:         pointer.To(model.FriendlyName),
					Description:          pointer.To(model.Description),
					HostPoolArmPath:      model.HostPoolId,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating Virtual Desktop Application Group %q (Resource Group %q): %+v", model.Name, model.ResourceGroupName, err)
			}

			if applicationgroup.ApplicationGroupType(model.Type) == applicationgroup.ApplicationGroupTypeDesktop {
				if model.DefaultDesktopDisplayName != "" {
					desktopClient := metadata.Client.DesktopVirtualization.DesktopsClient
					// default desktop name created for Application Group is 'sessionDesktop'
					desktopId := desktop.NewDesktopID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName, "sessionDesktop")
					desktopModel, err := desktopClient.Get(ctx, desktopId)
					if err != nil {
						if !response.WasNotFound(desktopModel.HttpResponse) {
							return fmt.Errorf("retrieving default desktop for %s: %+v", id, err)
						}
					}

					desktopPatch := desktop.DesktopPatch{
						Properties: &desktop.DesktopPatchProperties{
							FriendlyName: pointer.To(model.DefaultDesktopDisplayName),
						},
					}

					if _, err := desktopClient.Update(ctx, desktopId, desktopPatch); err != nil {
						return fmt.Errorf("setting friendly name for default desktop %s: %+v", id, err)
					}
				}
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (DesktopVirtualizationApplicationGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ApplicationGroupsClient

			state := DesktopVirtualizationApplicationGroupResourceModel{}

			id, err := applicationgroup.ParseApplicationGroupID(metadata.ResourceData.Id())
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

			state.Name = id.ApplicationGroupName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				props := model.Properties

				state.FriendlyName = pointer.From(props.FriendlyName)
				state.Description = pointer.From(props.Description)
				state.Type = string(props.ApplicationGroupType)
				defaultDesktopDisplayName := ""
				if props.ApplicationGroupType == applicationgroup.ApplicationGroupTypeDesktop {
					desktopClient := metadata.Client.DesktopVirtualization.DesktopsClient
					// default desktop name created for Application Group is 'sessionDesktop'
					desktopId := desktop.NewDesktopID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName, "sessionDesktop")
					desktopResp, err := desktopClient.Get(ctx, desktopId)
					if err != nil {
						if !response.WasNotFound(desktopResp.HttpResponse) {
							return fmt.Errorf("retrieving default desktop for %s: %+v", *id, err)
						}
					}
					// if the default desktop was found then set the display name attribute
					if desktopModel := desktopResp.Model; desktopModel != nil && desktopModel.Properties != nil && desktopModel.Properties.FriendlyName != nil {
						defaultDesktopDisplayName = *desktopModel.Properties.FriendlyName
					}
				}
				state.DefaultDesktopDisplayName = defaultDesktopDisplayName

				hostPoolId, err := hostpool.ParseHostPoolIDInsensitively(props.HostPoolArmPath)
				if err != nil {
					return fmt.Errorf("parsing Host Pool ID %q: %+v", props.HostPoolArmPath, err)
				}
				state.HostPoolId = hostPoolId.ID()
				state.Tags = pointer.From(model.Tags)
			}
			return metadata.Encode(&state)
		},
	}
}

func (r DesktopVirtualizationApplicationGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ApplicationGroupsClient

			var model DesktopVirtualizationApplicationGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			locks.ByName(model.Name, r.ResourceType())
			defer locks.UnlockByName(model.Name, r.ResourceType())

			id, err := applicationgroup.ParseApplicationGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("friendly_name") {
				payload.Properties.FriendlyName = pointer.To(model.FriendlyName)
			}
			if metadata.ResourceData.HasChange("default_desktop_display_name") &&
				applicationgroup.ApplicationGroupType(model.Type) == applicationgroup.ApplicationGroupTypeDesktop &&
				model.DefaultDesktopDisplayName != "" {
				desktopClient := metadata.Client.DesktopVirtualization.DesktopsClient
				// default desktop name created for Application Group is 'sessionDesktop'
				desktopId := desktop.NewDesktopID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName, "sessionDesktop")
				desktopModel, err := desktopClient.Get(ctx, desktopId)
				if err != nil {
					if !response.WasNotFound(desktopModel.HttpResponse) {
						return fmt.Errorf("retrieving default desktop for %s: %+v", id, err)
					}
				}

				desktopPatch := desktop.DesktopPatch{
					Properties: &desktop.DesktopPatchProperties{
						FriendlyName: pointer.To(model.DefaultDesktopDisplayName),
					},
				}

				if _, err := desktopClient.Update(ctx, desktopId, desktopPatch); err != nil {
					return fmt.Errorf("setting friendly name for default desktop %s: %+v", id, err)
				}
			}
			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}
			if _, err := client.CreateOrUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DesktopVirtualizationApplicationGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ApplicationGroupsClient

			id, err := applicationgroup.ParseApplicationGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.ApplicationGroupName, r.ResourceType())
			defer locks.UnlockByName(id.ApplicationGroupName, r.ResourceType())

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
