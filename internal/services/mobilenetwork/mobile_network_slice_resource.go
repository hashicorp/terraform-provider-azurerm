// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/slice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SliceResourceModel struct {
	Name                                             string                                                          `tfschema:"name"`
	MobileNetworkId                                  string                                                          `tfschema:"mobile_network_id"`
	Description                                      string                                                          `tfschema:"description"`
	Location                                         string                                                          `tfschema:"location"`
	SingleNetworkSliceSelectionAssistanceInformation []SingleNetworkSliceSelectionAssistanceInformationResourceModel `tfschema:"single_network_slice_selection_assistance_information"`
	Tags                                             map[string]string                                               `tfschema:"tags"`
}

type SingleNetworkSliceSelectionAssistanceInformationResourceModel struct {
	SliceDifferentiator string `tfschema:"slice_differentiator"`
	SliceServiceType    int64  `tfschema:"slice_service_type"`
}

type SliceResource struct{}

var _ sdk.ResourceWithUpdate = SliceResource{}

func (r SliceResource) ResourceType() string {
	return "azurerm_mobile_network_slice"
}

func (r SliceResource) ModelObject() interface{} {
	return &SliceResourceModel{}
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
					// TODO: these fields can be moved to the top-level in 4.0

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
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SliceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SliceClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(model.MobileNetworkId)
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

			properties := slice.Slice{
				Location:   location.Normalize(model.Location),
				Properties: slice.SlicePropertiesFormat{},
				Tags:       &model.Tags,
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			properties.Properties.Snssai = expandSingleNetworkSliceSelectionAssistanceInformationResourceModel(model.SingleNetworkSliceSelectionAssistanceInformation)

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SliceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SliceClient

			id, err := slice.ParseSliceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SliceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: properties were nil", id)
			}

			updateModel := resp.Model

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					updateModel.Properties.Description = &model.Description
				} else {
					updateModel.Properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("snssai") {
				updateModel.Properties.Snssai = expandSingleNetworkSliceSelectionAssistanceInformationResourceModel(model.SingleNetworkSliceSelectionAssistanceInformation)
			}

			if metadata.ResourceData.HasChange("tags") {
				updateModel.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *updateModel); err != nil {
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

			state := SliceResourceModel{
				Name:            id.SliceName,
				MobileNetworkId: mobilenetwork.NewMobileNetworkID(id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName).ID(),
			}
			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if model.Properties.Description != nil {
					state.Description = *model.Properties.Description
				}

				state.SingleNetworkSliceSelectionAssistanceInformation = flattenSingleNetworkSliceSelectionAssistanceInformationResourceModel(model.Properties.Snssai)

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SliceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
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

func expandSingleNetworkSliceSelectionAssistanceInformationResourceModel(input []SingleNetworkSliceSelectionAssistanceInformationResourceModel) slice.Snssai {
	item := input[0]

	output := slice.Snssai{
		Sst: item.SliceServiceType,
	}
	if item.SliceDifferentiator != "" {
		output.Sd = pointer.To(item.SliceDifferentiator)
	}

	return output
}

func flattenSingleNetworkSliceSelectionAssistanceInformationResourceModel(input slice.Snssai) []SingleNetworkSliceSelectionAssistanceInformationResourceModel {
	output := SingleNetworkSliceSelectionAssistanceInformationResourceModel{
		SliceServiceType: input.Sst,
	}

	if input.Sd != nil {
		output.SliceDifferentiator = *input.Sd
	}

	return []SingleNetworkSliceSelectionAssistanceInformationResourceModel{
		output,
	}
}
