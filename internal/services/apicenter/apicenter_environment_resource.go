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
	Name           string `tfschema:"name"`
	ServiceId      string `tfschema:"service_id"`
	Identification string `tfschema:"identification"`
	Description    string `tfschema:"description"`
	Type           string `tfschema:"environment_type"`
	DevPortalUri   string `tfschema:"development_portal_uri"`
	Instructions   string `tfschema:"instructions"`
	ServerType     string `tfschema:"server_type"`
	MgmtPortalUri  string `tfschema:"management_portal_uri"`
}

func (r ApiCenterEnvironmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"service_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"identification": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"development_portal_uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"instructions": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"server_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice(
				environments.PossibleValuesForEnvironmentServerType(),
				false),
		},
		"management_portal_uri": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
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
	return "azurerm_apicenter_environment"
}

func (r ApiCenterEnvironmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return environments.ValidateEnvironmentID
}

func (r ApiCenterEnvironmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApiCenterEnvironmentResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.ApiCenter.EnvironmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			service, err := services.ParseServiceID(model.ServiceId)
			if err != nil {
				return fmt.Errorf("parsing ApiCenter Service ID %s: %+v", model.ServiceId, err)
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
				Kind:        environments.EnvironmentKind(model.Type),
				Description: &model.Description,
				Onboarding: &environments.Onboarding{
					DeveloperPortalUri: &[]string{model.DevPortalUri},
					Instructions:       &model.Instructions,
				},
				Server: &environments.EnvironmentServer{
					Type:                pointer.To(environments.EnvironmentServerType(model.ServerType)),
					ManagementPortalUri: &[]string{model.MgmtPortalUri},
				},
			}

			apiCenterEnvironment := environments.Environment{
				Id:         &model.Identification,
				Name:       &model.Name,
				Properties: &apiCenterEnvironmentProps,
				Type:       &model.Type,
			}

			if _, err = client.CreateOrUpdate(ctx, id, apiCenterEnvironment); err != nil {
				return fmt.Errorf("creating ApiCenter Environment %s: %+v", id, err)
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
				return fmt.Errorf("reading ApiCenter Environment %s: %v", id, err)
			}

			if _, err = client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating ApiCenter Service %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ApiCenterEnvironmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := environments.ParseEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.ApiCenter.EnvironmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving ApiCenter Environment %s: %+v", *id, err)
			}

			serviceId := services.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)

			state := ApiCenterEnvironmentResourceModel{
				Name:           *resp.Model.Name,
				ServiceId:      serviceId.ID(),
				Identification: *resp.Model.Id,
				Description:    *resp.Model.Properties.Description,
				Type:           *resp.Model.Type,
				DevPortalUri:   (*resp.Model.Properties.Onboarding.DeveloperPortalUri)[0],
				MgmtPortalUri:  (*resp.Model.Properties.Server.ManagementPortalUri)[0],
				Instructions:   *resp.Model.Properties.Onboarding.Instructions,
				ServerType:     string(*resp.Model.Properties.Server.Type),
			}
			return metadata.Encode(&state)

		},
	}
}

func (r ApiCenterEnvironmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := environments.ParseEnvironmentID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.ApiCenter.EnvironmentsClient

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting ApiCenter Environment %s: %+v", *id, err)
			}

			return nil
		},
	}
}
