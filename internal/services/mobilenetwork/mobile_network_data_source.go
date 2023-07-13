// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MobileNetworkDataSource struct{}

type MobileNetworkDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	MobileCountryCode string            `tfschema:"mobile_country_code"`
	MobileNetworkCode string            `tfschema:"mobile_network_code"`
	Tags              map[string]string `tfschema:"tags"`
	ServiceKey        string            `tfschema:"service_key"`
}

var _ sdk.DataSource = MobileNetworkDataSource{}

func (r MobileNetworkDataSource) ResourceType() string {
	return "azurerm_mobile_network"
}

func (r MobileNetworkDataSource) ModelObject() interface{} {
	return &MobileNetworkDataSourceModel{}
}

func (r MobileNetworkDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return mobilenetwork.ValidateMobileNetworkID
}

func (r MobileNetworkDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r MobileNetworkDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"location": commonschema.LocationComputed(),

		"mobile_country_code": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"mobile_network_code": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),

		"service_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r MobileNetworkDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel MobileNetworkDataSourceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.MobileNetworkClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := mobilenetwork.NewMobileNetworkID(subscriptionId, metaModel.ResourceGroupName, metaModel.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := MobileNetworkDataSourceModel{
				Name:              id.MobileNetworkName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
				MobileCountryCode: model.Properties.PublicLandMobileNetworkIdentifier.Mcc,
				MobileNetworkCode: model.Properties.PublicLandMobileNetworkIdentifier.Mnc,
			}

			if model.Properties.ServiceKey != nil {
				state.ServiceKey = *model.Properties.ServiceKey
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
