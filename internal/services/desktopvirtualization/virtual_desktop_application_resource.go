// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/applicationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = DesktopVirtualizationApplicationResource{}
	_ sdk.ResourceWithUpdate = DesktopVirtualizationApplicationResource{}
)

type DesktopVirtualizationApplicationResource struct{}

func (DesktopVirtualizationApplicationResource) ModelObject() interface{} {
	return &DesktopVirtualizationApplicationResourceModel{}
}

func (DesktopVirtualizationApplicationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return application.ValidateApplicationID
}

func (DesktopVirtualizationApplicationResource) ResourceType() string {
	return "azurerm_virtual_desktop_application"
}

type DesktopVirtualizationApplicationResourceModel struct {
	Name                      string `tfschema:"name"`
	ApplicationGroupId        string `tfschema:"application_group_id"`
	FriendlyName              string `tfschema:"friendly_name"`
	Description               string `tfschema:"description"`
	Path                      string `tfschema:"path"`
	CommandLineArgumentPolicy string `tfschema:"command_line_argument_policy"`
	CommandLineArguments      string `tfschema:"command_line_arguments"`
	ShowInPortal              bool   `tfschema:"show_in_portal"`
	IconPath                  string `tfschema:"icon_path"`
	IconIndex                 int64  `tfschema:"icon_index"`
}

func (DesktopVirtualizationApplicationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{1,260}$"),
					"Virtual desktop application name must be 1 - 260 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},

		"application_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: applicationgroup.ValidateApplicationGroupID,
		},

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
			// NOTE: O+C The API will use the value in `name` as the default
			Computed: true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"path": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"command_line_argument_policy": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(application.PossibleValuesForCommandLineSetting(), false),
		},

		"command_line_arguments": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"show_in_portal": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"icon_path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C The API will use the value in `path` as the default
			Computed: true,
		},

		"icon_index": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
	}
}

func (DesktopVirtualizationApplicationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DesktopVirtualizationApplicationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ApplicationsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DesktopVirtualizationApplicationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			log.Printf("[INFO] preparing arguments for Virtual Desktop Application creation")

			applicationGroup, _ := applicationgroup.ParseApplicationGroupID(model.ApplicationGroupId)
			id := application.NewApplicationID(subscriptionId, applicationGroup.ResourceGroupName, applicationGroup.ApplicationGroupName, model.Name)

			locks.ByName(id.ApplicationName, r.ResourceType())
			defer locks.UnlockByName(id.ApplicationName, r.ResourceType())

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := application.Application{
				Properties: application.ApplicationProperties{
					FriendlyName:         pointer.To(model.FriendlyName),
					Description:          pointer.To(model.Description),
					FilePath:             pointer.To(model.Path),
					CommandLineSetting:   application.CommandLineSetting(model.CommandLineArgumentPolicy),
					CommandLineArguments: pointer.To(model.CommandLineArguments),
					ShowInPortal:         pointer.To(model.ShowInPortal),
					IconPath:             pointer.To(model.IconPath),
					IconIndex:            pointer.To(model.IconIndex),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DesktopVirtualizationApplicationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ApplicationsClient

			state := DesktopVirtualizationApplicationResourceModel{}

			id, err := application.ParseApplicationID(metadata.ResourceData.Id())
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

			state.Name = id.ApplicationName
			state.ApplicationGroupId = applicationgroup.NewApplicationGroupID(id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName).ID()

			if model := resp.Model; model != nil {
				props := model.Properties

				state.FriendlyName = pointer.From(props.FriendlyName)
				state.Description = pointer.From(props.Description)
				state.Path = pointer.From(props.FilePath)
				state.CommandLineArgumentPolicy = string(props.CommandLineSetting)
				state.CommandLineArguments = pointer.From(props.CommandLineArguments)
				state.ShowInPortal = pointer.From(props.ShowInPortal)
				state.IconPath = pointer.From(props.IconPath)
				state.IconIndex = pointer.From(props.IconIndex)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DesktopVirtualizationApplicationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ApplicationsClient

			var model DesktopVirtualizationApplicationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			locks.ByName(model.Name, r.ResourceType())
			defer locks.UnlockByName(model.Name, r.ResourceType())

			id, err := application.ParseApplicationID(metadata.ResourceData.Id())
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
			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}
			if metadata.ResourceData.HasChange("path") {
				payload.Properties.FilePath = pointer.To(model.Path)
			}
			if metadata.ResourceData.HasChange("command_line_argument_policy") {
				payload.Properties.CommandLineSetting = application.CommandLineSetting(model.CommandLineArgumentPolicy)
			}
			if metadata.ResourceData.HasChange("command_line_arguments") {
				payload.Properties.CommandLineArguments = pointer.To(model.CommandLineArguments)
			}
			if metadata.ResourceData.HasChange("show_in_portal") {
				payload.Properties.ShowInPortal = pointer.To(model.ShowInPortal)
			}
			if metadata.ResourceData.HasChange("icon_path") {
				payload.Properties.IconPath = pointer.To(model.IconPath)
			}
			if metadata.ResourceData.HasChange("icon_index") {
				payload.Properties.IconIndex = pointer.To(model.IconIndex)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DesktopVirtualizationApplicationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ApplicationsClient

			id, err := application.ParseApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.ApplicationName, r.ResourceType())
			defer locks.UnlockByName(id.ApplicationName, r.ResourceType())

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
