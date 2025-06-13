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

var _ sdk.DataSource = DevCenterEnvironmentTypeDataSource{}

type DevCenterEnvironmentTypeDataSource struct{}

type DevCenterEnvironmentTypeDataSourceModel struct {
	Name        string            `tfschema:"name"`
	DevCenterId string            `tfschema:"dev_center_id"`
	Tags        map[string]string `tfschema:"tags"`
}

func (DevCenterEnvironmentTypeDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DevCenterEnvironmentTypeName,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequired(&environmenttypes.DevCenterId{}),
	}
}

func (DevCenterEnvironmentTypeDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"tags": commonschema.TagsDataSource(),
	}
}

func (DevCenterEnvironmentTypeDataSource) ModelObject() interface{} {
	return &DevCenterEnvironmentTypeDataSourceModel{}
}

func (DevCenterEnvironmentTypeDataSource) ResourceType() string {
	return "azurerm_dev_center_environment_type"
}

func (r DevCenterEnvironmentTypeDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.EnvironmentTypes
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DevCenterEnvironmentTypeDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := devcenters.ParseDevCenterID(state.DevCenterId)
			if err != nil {
				return err
			}

			id := environmenttypes.NewDevCenterEnvironmentTypeID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, state.Name)

			resp, err := client.EnvironmentTypesGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Name = id.EnvironmentTypeName
				state.DevCenterId = environmenttypes.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID()
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}
