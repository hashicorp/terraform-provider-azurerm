// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apicenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiCenterServiceResource struct{}

var _ sdk.ResourceWithUpdate = ApiCenterServiceResource{}

type ApiCenterServiceResourceModel struct {
	Name          string            `tfschema:"name"`
	ResourceGroup string            `tfschema:"resource_group_name"`
	Location      string            `tfschema:"location"`
	Tags          map[string]string `tfschema:"tags"`
}

func (r ApiCenterServiceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccessConnectorName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (r ApiCenterServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiCenterServiceResource) ModelObject() interface{} {
	return &ApiCenterServiceResourceModel{}
}

func (r ApiCenterServiceResource) ResourceType() string {
	return "azurerm_apicenter_service"
}

func (r ApiCenterServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return services.ValidateServiceID
}

func (r ApiCenterServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApiCenterServiceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.ApiCenter.ServicesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := services.NewServiceID(subscriptionId, model.ResourceGroup, model.Name)
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

func (r ApiCenterServiceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiCenter.ServicesClient
			id, err := services.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ApiCenterServiceResourceModel
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

func (r ApiCenterServiceResource) Read() sdk.ResourceFunc {
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

			state := ApiCenterServiceResourceModel{
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

func (r ApiCenterServiceResource) Delete() sdk.ResourceFunc {
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
