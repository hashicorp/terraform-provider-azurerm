package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/slice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SliceDataSource struct{}

var _ sdk.DataSource = SliceDataSource{}

func (r SliceDataSource) ResourceType() string {
	return "azurerm_mobile_network_slice"
}

func (r SliceDataSource) ModelObject() interface{} {
	return &SliceModel{}
}

func (r SliceDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return slice.ValidateSliceID
}

func (r SliceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: mobilenetwork.ValidateMobileNetworkID,
		},
	}
}

func (r SliceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"single_network_slice_selection_assistance_information": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"slice_differentiator": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"slice_service_type": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r SliceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel SliceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SliceClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(metaModel.MobileNetworkId)
			if err != nil {
				return err
			}

			id := slice.NewSliceID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, metaModel.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			model := *resp.Model

			state := SliceModel{
				Name:            id.SliceName,
				MobileNetworkId: mobileNetworkId.ID(),
				Location:        location.Normalize(model.Location),
			}

			properties := model.Properties
			if properties.Description != nil {
				state.Description = *properties.Description
			}

			state.Snssai = flattenSnssaiModel(properties.Snssai)

			if resp.Model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
