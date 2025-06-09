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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DevCenterDevBoxDefinitionDataSource{}

type DevCenterDevBoxDefinitionDataSource struct{}

type DevCenterDevBoxDefinitionDataSourceModel struct {
	Name             string            `tfschema:"name"`
	Location         string            `tfschema:"location"`
	DevCenterId      string            `tfschema:"dev_center_id"`
	ImageReferenceId string            `tfschema:"image_reference_id"`
	SkuName          string            `tfschema:"sku_name"`
	Tags             map[string]string `tfschema:"tags"`
}

func (DevCenterDevBoxDefinitionDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DevCenterDevBoxDefinitionName,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequired(&devboxdefinitions.DevCenterId{}),
	}
}

func (DevCenterDevBoxDefinitionDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"image_reference_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (DevCenterDevBoxDefinitionDataSource) ModelObject() interface{} {
	return &DevCenterDevBoxDefinitionDataSourceModel{}
}

func (DevCenterDevBoxDefinitionDataSource) ResourceType() string {
	return "azurerm_dev_center_dev_box_definition"
}

func (r DevCenterDevBoxDefinitionDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevBoxDefinitions
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DevCenterDevBoxDefinitionDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := devboxdefinitions.ParseDevCenterID(state.DevCenterId)
			if err != nil {
				return err
			}

			id := devboxdefinitions.NewDevCenterDevBoxDefinitionID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state.Name = id.DevBoxDefinitionName
			state.DevCenterId = devboxdefinitions.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID()

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					if v := props.ImageReference; v != nil {
						state.ImageReferenceId = pointer.From(v.Id)
					}

					if v := props.Sku; v != nil {
						state.SkuName = v.Name
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
