package mobilenetwork

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/slice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SliceModel struct {
	Name                         string            `tfschema:"name"`
	MobileNetworkMobileNetworkId string            `tfschema:"mobile_network_id"`
	Description                  string            `tfschema:"description"`
	Location                     string            `tfschema:"location"`
	Snssai                       []SnssaiModel     `tfschema:"single_network_slice_selection_assistance_information"`
	Tags                         map[string]string `tfschema:"tags"`
}

type SnssaiModel struct {
	Sd  string `tfschema:"slice_differentiator"`
	Sst int64  `tfschema:"slice_service_type"`
}

type SliceResource struct{}

var _ sdk.ResourceWithUpdate = SliceResource{}

func (r SliceResource) ResourceType() string {
	return "azurerm_mobile_network_slice"
}

func (r SliceResource) ModelObject() interface{} {
	return &SliceModel{}
}

func (r SliceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return slice.ValidateSliceID
}

func (r SliceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: mobilenetwork.ValidateMobileNetworkID,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"single_network_slice_selection_assistance_information": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"slice_differentiator": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[A-Fa-f0-9]{6}$`),
							"Slice Differentiator must be a 6 digit hex string",
						),
					},

					"slice_service_type": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 255),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r SliceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SliceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SliceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SliceClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(model.MobileNetworkMobileNetworkId)
			if err != nil {
				return err
			}

			id := slice.NewSliceID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &slice.Slice{
				Location:   location.Normalize(model.Location),
				Properties: slice.SlicePropertiesFormat{},
				Tags:       &model.Tags,
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			snssaiValue, err := expandSnssaiModel(model.Snssai)
			if err != nil {
				return err
			}

			if snssaiValue != nil {
				properties.Properties.Snssai = *snssaiValue
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SliceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SliceClient

			id, err := slice.ParseSliceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SliceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Properties.Description = &model.Description
				} else {
					properties.Properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("snssai") {
				snssaiValue, err := expandSnssaiModel(model.Snssai)
				if err != nil {
					return err
				}

				if snssaiValue != nil {
					properties.Properties.Snssai = *snssaiValue
				}
			}

			properties.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SliceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SliceClient

			id, err := slice.ParseSliceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := SliceModel{
				Name:                         id.SliceName,
				MobileNetworkMobileNetworkId: mobilenetwork.NewMobileNetworkID(id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName).ID(),
				Location:                     location.Normalize(model.Location),
			}

			properties := &model.Properties
			if properties.Description != nil {
				state.Description = *properties.Description
			}

			snssaiValue, err := flattenSnssaiModel(&properties.Snssai)
			if err != nil {
				return err
			}

			state.Snssai = snssaiValue
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SliceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SliceClient

			id, err := slice.ParseSliceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandSnssaiModel(inputList []SnssaiModel) (*slice.Snssai, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := slice.Snssai{
		Sst: input.Sst,
	}

	if input.Sd != "" {
		output.Sd = &input.Sd
	}

	return &output, nil
}

func flattenSnssaiModel(input *slice.Snssai) ([]SnssaiModel, error) {
	var outputList []SnssaiModel
	if input == nil {
		return outputList, nil
	}

	output := SnssaiModel{
		Sst: input.Sst,
	}

	if input.Sd != nil {
		output.Sd = *input.Sd
	}

	return append(outputList, output), nil
}
