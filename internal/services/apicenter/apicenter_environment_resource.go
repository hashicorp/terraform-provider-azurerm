// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apicenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
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

			id := environments.NewEnvironmentID(subscriptionId, model.ResourceGroup, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing ApiCenter Service %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			apiCenterService := services.Service{
				Name:     &model.Name,
				Location: model.Location,
				Tags:     &model.Tags,
				Identity: expandedIdentity,
			}

			if _, err = client.CreateOrUpdate(ctx, id, apiCenterService); err != nil {
				return fmt.Errorf("creating ApiCenter Service %s: %+v", id, err)
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
			client := metadata.Client.ApiCenter.ServicesClient
			id, err := services.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ApiCenterEnvironmentResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading ApiCenter Service %s: %v", id, err)
			}

			if metadata.ResourceData.HasChange("identity") {
				// TODO: Switch this to 'identity.ExpandSystemOrSingleUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))'
				// once SDK Helpers PR #164 has been merged and integrated into the provider...
				identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				existing.Model.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = &state.Tags
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
			id, err := services.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.ApiCenter.ServicesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving ApiCenter Service %s: %+v", *id, err)
			}

			state := ApiCenterEnvironmentResourceModel{
				Name:          id.ServiceName,
				Location:      location.Normalize(resp.Model.Location),
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if model.Identity != nil {
					identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
					if err != nil {
						return fmt.Errorf("flattening `identity`: %+v", err)
					}

					if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
						return fmt.Errorf("setting `identity`: %+v", err)
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
			id, err := services.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.ApiCenter.ServicesClient

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting ApiCenter Service %s: %+v", *id, err)
			}

			return nil
		},
	}
}
