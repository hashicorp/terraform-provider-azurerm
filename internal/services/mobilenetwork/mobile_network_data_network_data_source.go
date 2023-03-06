package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/datanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataNetworkDataSource struct{}

var _ sdk.DataSource = DataNetworkDataSource{}

func (r DataNetworkDataSource) ResourceType() string {
	return "azurerm_mobile_network_data_network"
}

func (r DataNetworkDataSource) ModelObject() interface{} {
	return &DataNetworkModel{}
}

func (r DataNetworkDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datanetwork.ValidateDataNetworkID
}

func (r DataNetworkDataSource) Arguments() map[string]*pluginsdk.Schema {
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

func (r DataNetworkDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DataNetworkDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var inputModel DataNetworkModel
			if err := metadata.Decode(&inputModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.DataNetworkClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(inputModel.MobileNetworkMobileNetworkId)
			if err != nil {
				return err
			}

			id := datanetwork.NewDataNetworkID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, inputModel.Name)

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

			state := DataNetworkModel{
				Name:                         id.DataNetworkName,
				MobileNetworkMobileNetworkId: mobilenetwork.NewMobileNetworkID(id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName).ID(),
				Location:                     location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				if properties.Description != nil {
					state.Description = *properties.Description
				}
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
