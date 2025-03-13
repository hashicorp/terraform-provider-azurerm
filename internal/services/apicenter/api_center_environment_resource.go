// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apicenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/environments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiCenterEnvironmentResource struct{}

var _ sdk.ResourceWithUpdate = ApiCenterEnvironmentResource{}

type ApiCenterEnvironmentResourceModel struct {
	Name          string `tfschema:"name"`
	ServiceId     string `tfschema:"api_center_service_id"`
	Title         string `tfschema:"title"`
	Description   string `tfschema:"description"`
	Type          string `tfschema:"environment_type"`
	DevPortalUri  string `tfschema:"development_portal_uri"`
	Instructions  string `tfschema:"instructions"`
	ServerType    string `tfschema:"server_type"`
	MgmtPortalUri string `tfschema:"management_portal_uri"`
}

func (r ApiCenterEnvironmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"api_center_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: services.ValidateServiceID,
		},

		"title": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"environment_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				environments.PossibleValuesForEnvironmentKind(),
				false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"development_portal_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"instructions": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"server_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice(
				environments.PossibleValuesForEnvironmentServerType(),
				false),
		},
		"management_portal_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},
	}
}

func (r ApiCenterEnvironmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiCenterEnvironmentResource) ModelObject() interface{} {
	return &ApiCenterEnvironmentResourceModel{}
}

func (r ApiCenterEnvironmentResource) ResourceType() string {
	return "azurerm_api_center_environment"
}

func (r ApiCenterEnvironmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return environments.ValidateEnvironmentID
}

func (r ApiCenterEnvironmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.EnvironmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiCenterEnvironmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			service, err := services.ParseServiceID(model.ServiceId)
			if err != nil {
				return err
			}

			// @favoretti: can't find any workspace creation buttons in the portal, assume everything defaults to "default" for now?
			id := environments.NewEnvironmentID(subscriptionId, service.ResourceGroupName, service.ServiceName, "default", model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing ApiCenter Environment %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			apiCenterEnvironmentProps := environments.EnvironmentProperties{
				Kind:  environments.EnvironmentKind(model.Type),
				Title: model.Title,
			}

			if model.Description != "" {
				apiCenterEnvironmentProps.Description = pointer.To(model.Description)
			}

			apiCenterEnvironmentProps.Onboarding = &environments.Onboarding{}
			if model.DevPortalUri != "" {
				apiCenterEnvironmentProps.Onboarding.DeveloperPortalUri = pointer.To([]string{model.DevPortalUri})
			}

			if model.Instructions != "" {
				apiCenterEnvironmentProps.Onboarding.Instructions = pointer.To(model.Instructions)
			}

			apiCenterEnvironmentProps.Server = &environments.EnvironmentServer{}

			if model.ServerType != "" {
				apiCenterEnvironmentProps.Server.Type = pointer.To(environments.EnvironmentServerType(model.ServerType))
			}

			if model.MgmtPortalUri != "" {
				apiCenterEnvironmentProps.Server.ManagementPortalUri = pointer.To([]string{model.MgmtPortalUri})
			}

			apiCenterEnvironment := environments.Environment{
				Name:       pointer.To(model.Name),
				Properties: pointer.To(apiCenterEnvironmentProps),
				Type:       pointer.To(model.Type),
			}

			if _, err = client.CreateOrUpdate(ctx, id, apiCenterEnvironment); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiCenterEnvironmentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.EnvironmentsClient
			id, err := environments.ParseEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ApiCenterEnvironmentResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			model := existing.Model

			if metadata.ResourceData.HasChange("description") {
				model.Properties.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("development_portal_uri") {
				model.Properties.Onboarding.DeveloperPortalUri = pointer.To([]string{state.DevPortalUri})
			}

			if metadata.ResourceData.HasChange("server_type") {
				model.Properties.Server.Type = pointer.To(environments.EnvironmentServerType(state.ServerType))
			}

			if metadata.ResourceData.HasChange("management_portal_uri") {
				model.Properties.Server.ManagementPortalUri = pointer.To([]string{state.MgmtPortalUri})
			}

			if metadata.ResourceData.HasChange("instructions") {
				model.Properties.Onboarding.Instructions = pointer.To(state.Instructions)
			}

			if _, err = client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ApiCenterEnvironmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.EnvironmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id, err := environments.ParseEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiCenterEnvironmentResourceModel{
				ServiceId: services.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName).ID(),
			}
			if model := existing.Model; model != nil {
				state.Name = pointer.From(model.Name)
				if props := existing.Model.Properties; props != nil {
					state.Title = props.Title
					state.Description = pointer.From(props.Description)
					state.Type = string(props.Kind)
					if server := props.Server; server != nil {
						if pointer.From(server.ManagementPortalUri) != nil && len(pointer.From(server.ManagementPortalUri)) != 0 {
							state.MgmtPortalUri = pointer.From(server.ManagementPortalUri)[0]
						}
						state.ServerType = string(pointer.From(server.Type))
					}

					if onboarding := props.Onboarding; onboarding != nil {
						if pointer.From(onboarding.DeveloperPortalUri) != nil && len(pointer.From(onboarding.DeveloperPortalUri)) != 0 {
							state.DevPortalUri = pointer.From(onboarding.DeveloperPortalUri)[0]
						}
						state.Instructions = pointer.From(onboarding.Instructions)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiCenterEnvironmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.EnvironmentsClient

			id, err := environments.ParseEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
