// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = DevCenterEnvironmentTypeResource{}
	_ sdk.ResourceWithUpdate = DevCenterEnvironmentTypeResource{}
)

type DevCenterEnvironmentTypeResource struct{}

func (r DevCenterEnvironmentTypeResource) ModelObject() interface{} {
	return &DevCenterEnvironmentTypeResourceModel{}
}

type DevCenterEnvironmentTypeResourceModel struct {
	Name        string            `tfschema:"name"`
	DevCenterId string            `tfschema:"dev_center_id"`
	Tags        map[string]string `tfschema:"tags"`
}

func (r DevCenterEnvironmentTypeResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return environmenttypes.ValidateDevCenterEnvironmentTypeID
}

func (r DevCenterEnvironmentTypeResource) ResourceType() string {
	return "azurerm_dev_center_environment_type"
}

func (r DevCenterEnvironmentTypeResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterEnvironmentTypeName,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequiredForceNew(&environmenttypes.DevCenterId{}),

		"tags": commonschema.Tags(),
	}
}

func (r DevCenterEnvironmentTypeResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DevCenterEnvironmentTypeResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterEnvironmentTypeResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := devcenters.ParseDevCenterID(model.DevCenterId)
			if err != nil {
				return err
			}

			id := environmenttypes.NewDevCenterEnvironmentTypeID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, model.Name)

			existing, err := client.EnvironmentTypesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := environmenttypes.EnvironmentType{
				Tags: pointer.To(model.Tags),
			}

			if _, err := client.EnvironmentTypesCreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterEnvironmentTypeResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes

			id, err := environmenttypes.ParseDevCenterEnvironmentTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.EnvironmentTypesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterEnvironmentTypeResourceModel{
				Name:        id.EnvironmentTypeName,
				DevCenterId: environmenttypes.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DevCenterEnvironmentTypeResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes

			id, err := environmenttypes.ParseDevCenterEnvironmentTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.EnvironmentTypesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterEnvironmentTypeResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes

			id, err := environmenttypes.ParseDevCenterEnvironmentTypeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DevCenterEnvironmentTypeResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := environmenttypes.EnvironmentTypeUpdate{
				Tags: pointer.To(model.Tags),
			}

			if _, err := client.EnvironmentTypesUpdate(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
