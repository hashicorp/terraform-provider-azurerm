// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dev_center_dev_box_definition

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	devcenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

type DevCenterDevBoxDefinitionDataSourceModel struct {
	Name             string            `tfschema:"name"`
	Location         string            `tfschema:"location"`
	DevCenterId      string            `tfschema:"dev_center_id"`
	ImageReferenceId string            `tfschema:"image_reference_id"`
	HibernateSupport bool              `tfschema:"hibernate_support_enabled"`
	SkuName          string            `tfschema:"sku_name"`
	Tags             map[string]string `tfschema:"tags"`
}

func (DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Dev Center Dev Box Definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DevCenterDevBoxDefinitionName,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequired(&devboxdefinitions.DevCenterId{}, "Dev Center"),
	}
}

func (DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"image_reference_id": {
			Description: "The ID of the image for the Dev Center Dev Box Definition.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"location": commonschema.LocationComputed(),

		"hibernate_support_enabled": {
			Description: "Whether Dev Boxes created with this definition are capable of hibernation.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"sku_name": {
			Description: "The name of the SKU for the Dev Center Dev Box Definition.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (DataSource) ModelObject() any {
	return &DevCenterDevBoxDefinitionDataSourceModel{}
}

func (DataSource) ResourceType() string {
	return "azurerm_dev_center_dev_box_definition"
}

func (DataSource) Description() string {
	return "Use this data source to access information about an existing Dev Center Dev Box Definition."
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := devcenterClient.From(metadata.Client).V20250201.DevBoxDefinitions
			subscriptionId := account.From(metadata.Client).SubscriptionId

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

					state.HibernateSupport = pointer.From(props.HibernateSupport) == devboxdefinitions.HibernateSupportEnabled

					if v := props.Sku; v != nil {
						state.SkuName = v.Name
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
