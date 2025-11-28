// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SliceResourceModel struct {
	Name                                             string                                                          `tfschema:"name"`
	MobileNetworkId                                  string                                                          `tfschema:"mobile_network_id"`
	Description                                      string                                                          `tfschema:"description"`
	Location                                         string                                                          `tfschema:"location"`
	SliceDifferentiator                              string                                                          `tfschema:"slice_differentiator"`
	SliceServiceType                                 int64                                                           `tfschema:"slice_service_type"`
	SingleNetworkSliceSelectionAssistanceInformation []SingleNetworkSliceSelectionAssistanceInformationResourceModel `tfschema:"single_network_slice_selection_assistance_information,removedInNextMajorVersion"`
	Tags                                             map[string]string                                               `tfschema:"tags"`
}

type SingleNetworkSliceSelectionAssistanceInformationResourceModel struct {
	SliceDifferentiator string `tfschema:"slice_differentiator,removedInNextMajorVersion"`
	SliceServiceType    int64  `tfschema:"slice_service_type,removedInNextMajorVersion"`
}

type SliceResource struct{}

var _ sdk.ResourceWithUpdate = SliceResource{}

var _ sdk.ResourceWithDeprecationAndNoReplacement = SliceResource{}

func (r SliceResource) DeprecationMessage() string {
	return "The `azurerm_mobile_network_slice` resource has been deprecated and will be removed in v5.0 of the AzureRM Provider"
}

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
	s := map[string]*pluginsdk.Schema{
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

		"slice_service_type": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 255),
		},

		"slice_differentiator": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Fa-f0-9]{6}$`),
				"Slice Differentiator must be a 6 digit hex string",
			),
		},

		"tags": commonschema.Tags(),
	}

	if !features.FivePointOh() {
		s["slice_service_type"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			Computed:      true,
			ValidateFunc:  validation.IntBetween(0, 255),
			ConflictsWith: []string{"single_network_slice_selection_assistance_information"},
		}

		s["slice_differentiator"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Fa-f0-9]{6}$`),
				"Slice Differentiator must be a 6 digit hex string",
			),
			ConflictsWith: []string{"single_network_slice_selection_assistance_information"},
		}

		s["single_network_slice_selection_assistance_information"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeList,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"slice_service_type", "slice_differentiator"},
			MaxItems:      1,
			Deprecated:    "`single_network_slice_selection_assistance_information` has been deprecated and its properties, `slice_differentiator` and `slice_service_type` have been moved to the top level. The `single_network_slice_selection_assistance_information` block will be removed in v5.0 of the AzureRM Provider.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"slice_differentiator": {
						Type:       pluginsdk.TypeString,
						Optional:   true,
						Deprecated: "`single_network_slice_selection_assistance_information` has been deprecated and its properties, `slice_differentiator` and `slice_service_type` have been moved to the top level. The `single_network_slice_selection_assistance_information` block will be removed in v5.0 of the AzureRM Provider.",
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[A-Fa-f0-9]{6}$`),
							"Slice Differentiator must be a 6 digit hex string",
						),
					},

					"slice_service_type": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						Deprecated:   "`single_network_slice_selection_assistance_information` has been deprecated and its properties, `slice_differentiator` and `slice_service_type` have been moved to the top level. The `single_network_slice_selection_assistance_information` block will be removed in v5.0 of the AzureRM Provider.",
						ValidateFunc: validation.IntBetween(0, 255),
					},
				},
			},
		}
	}

	return s
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
				Location: location.Normalize(model.Location),
				Properties: slice.SlicePropertiesFormat{
					Snssai: slice.Snssai{
						Sst: model.SliceServiceType,
					},
				},
				Tags: &model.Tags,
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if model.SliceDifferentiator != "" {
				properties.Properties.Snssai.Sd = &model.SliceDifferentiator
			}

			if !features.FivePointOh() && len(model.SingleNetworkSliceSelectionAssistanceInformation) > 0 {
				properties.Properties.Snssai = expandSingleNetworkSliceSelectionAssistanceInformationResourceModel(model.SingleNetworkSliceSelectionAssistanceInformation)
			}

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

			if metadata.ResourceData.HasChange("single_network_slice_selection_assistance_information") {
				if !features.FivePointOh() {
					updateModel.Properties.Snssai = expandSingleNetworkSliceSelectionAssistanceInformationResourceModel(model.SingleNetworkSliceSelectionAssistanceInformation)
				}
			}
			if metadata.ResourceData.HasChange("single_network_slice_selection_assistance_information") {
				updateModel.Properties.Snssai = expandSingleNetworkSliceSelectionAssistanceInformationResourceModel(model.SingleNetworkSliceSelectionAssistanceInformation)
			}

			if metadata.ResourceData.HasChange("slice_service_type") {
				updateModel.Properties.Snssai.Sst = model.SliceServiceType
			}

			if metadata.ResourceData.HasChange("slice_differentiator") {
				updateModel.Properties.Snssai.Sd = &model.SliceDifferentiator
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

				if !features.FivePointOh() {
					state.SingleNetworkSliceSelectionAssistanceInformation = flattenSingleNetworkSliceSelectionAssistanceInformationResourceModel(model.Properties.Snssai)
				}
				state.SliceServiceType = model.Properties.Snssai.Sst
				if model.Properties.Snssai.Sd != nil {
					state.SliceDifferentiator = *model.Properties.Snssai.Sd
				}
				state.Tags = pointer.From(model.Tags)
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
