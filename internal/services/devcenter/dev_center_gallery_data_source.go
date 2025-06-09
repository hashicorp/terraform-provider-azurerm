// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = DevCenterGalleryDataSource{}

type DevCenterGalleryDataSource struct{}

type DevCenterGalleryDataSourceModel struct {
	DevCenterId       string `tfschema:"dev_center_id"`
	GalleryResourceId string `tfschema:"shared_gallery_id"`
	Name              string `tfschema:"name"`
}

func (DevCenterGalleryDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequired(&devcenters.DevCenterId{}),
	}
}

func (DevCenterGalleryDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"shared_gallery_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (DevCenterGalleryDataSource) ModelObject() interface{} {
	return &DevCenterGalleryDataSourceModel{}
}

func (DevCenterGalleryDataSource) ResourceType() string {
	return "azurerm_dev_center_gallery"
}

func (r DevCenterGalleryDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.Galleries
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DevCenterGalleryDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			devCenterId, err := galleries.ParseDevCenterID(state.DevCenterId)
			if err != nil {
				return err
			}

			id := galleries.NewGalleryID(subscriptionId, devCenterId.ResourceGroupName, devCenterId.DevCenterName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state.Name = id.GalleryName
			state.DevCenterId = galleries.NewDevCenterID(id.SubscriptionId, id.ResourceGroupName, id.DevCenterName).ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.GalleryResourceId = props.GalleryResourceId
				}
			}

			return metadata.Encode(&state)
		},
	}
}
