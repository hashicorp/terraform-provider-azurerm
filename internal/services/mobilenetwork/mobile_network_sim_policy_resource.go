// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/datanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/service"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/slice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SimPolicyResource struct{}

type SimPolicyResourceModel struct {
	Name                string                            `tfschema:"name"`
	MobileNetworkId     string                            `tfschema:"mobile_network_id"`
	DefaultSliceId      string                            `tfschema:"default_slice_id"`
	Location            string                            `tfschema:"location"`
	RegistrationTimer   int64                             `tfschema:"registration_timer_in_seconds"`
	RfspIndex           int64                             `tfschema:"rat_frequency_selection_priority_index"`
	SliceConfigurations []SliceConfigurationResourceModel `tfschema:"slice"`
	Tags                map[string]string                 `tfschema:"tags"`
	UeAmbr              []AmbrResourceModel               `tfschema:"user_equipment_aggregate_maximum_bit_rate"`
}

type SliceConfigurationResourceModel struct {
	DataNetworkConfigurations []DataNetworkConfigurationResourceModel `tfschema:"data_network"`
	DefaultDataNetworkId      string                                  `tfschema:"default_data_network_id"`
	SliceId                   string                                  `tfschema:"slice_id"`
}

type DataNetworkConfigurationResourceModel struct {
	AdditionalAllowedSessionTypes       []string            `tfschema:"additional_allowed_session_types"`
	AllocationAndRetentionPriorityLevel int64               `tfschema:"allocation_and_retention_priority_level"`
	AllowedServices                     []string            `tfschema:"allowed_services_ids"`
	DataNetworkId                       string              `tfschema:"data_network_id"`
	DefaultSessionType                  string              `tfschema:"default_session_type"`
	QosIdentifier                       int64               `tfschema:"qos_indicator"`
	MaximumNumberOfBufferedPackets      int64               `tfschema:"max_buffered_packets"`
	PreemptionCapability                string              `tfschema:"preemption_capability"`
	PreemptionVulnerability             string              `tfschema:"preemption_vulnerability"`
	SessionAmbr                         []AmbrResourceModel `tfschema:"session_aggregate_maximum_bit_rate"`
}

type AmbrResourceModel struct {
	Downlink string `tfschema:"downlink"`
	Uplink   string `tfschema:"uplink"`
}

var _ sdk.ResourceWithUpdate = SimPolicyResource{}

func (r SimPolicyResource) ResourceType() string {
	return "azurerm_mobile_network_sim_policy"
}

func (r SimPolicyResource) ModelObject() interface{} {
	return &SimPolicyResourceModel{}
}

func (r SimPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return simpolicy.ValidateSimPolicyID
}

func (r SimPolicyResource) Arguments() map[string]*pluginsdk.Schema {
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

		"default_slice_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: slice.ValidateSliceID,
		},

		"location": commonschema.Location(),

		"registration_timer_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      3240,
			ValidateFunc: validation.IntAtLeast(30),
		},

		"rat_frequency_selection_priority_index": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 256),
		},

		"slice": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_network": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"additional_allowed_session_types": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											string(simpolicy.PduSessionTypeIPvFour),
											string(simpolicy.PduSessionTypeIPvSix),
										}, false),
									},
								},

								"allocation_and_retention_priority_level": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntAtLeast(0),
								},

								"allowed_services_ids": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: service.ValidateServiceID,
									},
								},

								"data_network_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: datanetwork.ValidateDataNetworkID,
								},

								"default_session_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  simpolicy.PduSessionTypeIPvFour,
									ValidateFunc: validation.StringInSlice([]string{
										string(simpolicy.PduSessionTypeIPvFour),
										string(simpolicy.PduSessionTypeIPvSix),
									}, false),
								},

								"qos_indicator": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 127),
								},
								"max_buffered_packets": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      10,
									ValidateFunc: validation.IntAtLeast(0),
								},

								"preemption_capability": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  simpolicy.PreemptionCapabilityNotPreempt,
									ValidateFunc: validation.StringInSlice([]string{
										string(simpolicy.PreemptionCapabilityNotPreempt),
										string(simpolicy.PreemptionCapabilityMayPreempt),
									}, false),
								},

								"preemption_vulnerability": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  simpolicy.PreemptionVulnerabilityNotPreemptable,
									ValidateFunc: validation.StringInSlice([]string{
										string(simpolicy.PreemptionVulnerabilityNotPreemptable),
										string(simpolicy.PreemptionVulnerabilityPreemptable),
									}, false),
								},

								"session_aggregate_maximum_bit_rate": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"downlink": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringMatch(
													regexp.MustCompile(`^([1-9]\d*|0)(\.\d+)?\s(Kbps|Mbps|Gbps|Tbps)$`),
													"The value must be a number followed by Kbps, Mbps, Gbps or Tbps.",
												),
											},

											"uplink": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringMatch(
													regexp.MustCompile(`^([1-9]\d*|0)(\.\d+)?\s(Kbps|Mbps|Gbps|Tbps)$`),
													"The value must be a number followed by Kbps, Mbps, Gbps or Tbps.",
												),
											},
										},
									},
								},
							},
						},
					},

					"default_data_network_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: datanetwork.ValidateDataNetworkID,
					},

					"slice_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: slice.ValidateSliceID,
					},
				},
			},
		},

		"user_equipment_aggregate_maximum_bit_rate": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"downlink": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^([1-9]\d*|0)(\.\d+)?\s(Kbps|Mbps|Gbps|Tbps)$`),
							"The value must be a number followed by Kbps, Mbps, Gbps or Tbps.",
						),
					},

					"uplink": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^([1-9]\d*|0)(\.\d+)?\s(Kbps|Mbps|Gbps|Tbps)$`),
							"The value must be a number followed by Kbps, Mbps, Gbps or Tbps.",
						),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r SimPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SimPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SimPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SIMPolicyClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(model.MobileNetworkId)
			if err != nil {
				return err
			}

			id := simpolicy.NewSimPolicyID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, model.Name)
			existing, err := client.SimPoliciesGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := simpolicy.SimPolicy{
				Location: location.Normalize(model.Location),
				Properties: simpolicy.SimPolicyPropertiesFormat{
					RegistrationTimer: &model.RegistrationTimer,
					DefaultSlice: simpolicy.SliceResourceId{
						Id: model.DefaultSliceId,
					},
				},
				Tags: &model.Tags,
			}

			if model.RfspIndex != 0 {
				properties.Properties.RfspIndex = &model.RfspIndex
			}

			properties.Properties.SliceConfigurations = expandSliceConfigurationResourceModel(model.SliceConfigurations)

			properties.Properties.UeAmbr = expandAmbrResourceModel(model.UeAmbr)

			if err := client.SimPoliciesCreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r SimPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SIMPolicyClient

			id, err := simpolicy.ParseSimPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan SimPolicyResourceModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.SimPoliciesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			model := *resp.Model

			if metadata.ResourceData.HasChange("default_slice") {
				model.Properties.DefaultSlice = simpolicy.SliceResourceId{Id: plan.DefaultSliceId}

			}

			if metadata.ResourceData.HasChange("registration_timer_in_seconds") {
				model.Properties.RegistrationTimer = &plan.RegistrationTimer
			}

			if metadata.ResourceData.HasChange("rat_frequency_selection_priority_index") {
				model.Properties.RfspIndex = &plan.RfspIndex
			}

			if metadata.ResourceData.HasChange("slice") {
				model.Properties.SliceConfigurations = expandSliceConfigurationResourceModel(plan.SliceConfigurations)
			}

			if metadata.ResourceData.HasChange("user_equipment_aggregate_maximum_bit_rate") {
				model.Properties.UeAmbr = expandAmbrResourceModel(plan.UeAmbr)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = &plan.Tags
			}

			if err := client.SimPoliciesCreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SimPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SIMPolicyClient

			id, err := simpolicy.ParseSimPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.SimPoliciesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := SimPolicyResourceModel{
				Name:            id.SimPolicyName,
				MobileNetworkId: mobilenetwork.NewMobileNetworkID(id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				state.DefaultSliceId = model.Properties.DefaultSlice.Id

				if model.Properties.RegistrationTimer != nil {
					state.RegistrationTimer = *model.Properties.RegistrationTimer
				}

				if model.Properties.RfspIndex != nil {
					state.RfspIndex = *model.Properties.RfspIndex
				}

				state.SliceConfigurations = flattenSliceConfigurationResourceModel(model.Properties.SliceConfigurations)
				state.UeAmbr = flattenSimPolicyAmbrResourceModel(model.Properties.UeAmbr)

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SimPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SIMPolicyClient

			id, err := simpolicy.ParseSimPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.SimPoliciesDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandSliceConfigurationResourceModel(inputList []SliceConfigurationResourceModel) []simpolicy.SliceConfiguration {
	var outputList []simpolicy.SliceConfiguration
	for _, v := range inputList {
		input := v
		output := simpolicy.SliceConfiguration{
			Slice: simpolicy.SliceResourceId{
				Id: input.SliceId,
			},
		}

		output.DataNetworkConfigurations = expandDataNetworkConfigurationResourceModel(input.DataNetworkConfigurations)

		if input.DefaultDataNetworkId != "" {
			output.DefaultDataNetwork = simpolicy.DataNetworkResourceId{
				Id: input.DefaultDataNetworkId,
			}
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func expandDataNetworkConfigurationResourceModel(inputList []DataNetworkConfigurationResourceModel) []simpolicy.DataNetworkConfiguration {
	var outputList []simpolicy.DataNetworkConfiguration
	for _, v := range inputList {
		input := v

		defaultSessionType := simpolicy.PduSessionType(input.DefaultSessionType)
		preemptionCapability := simpolicy.PreemptionCapability(input.PreemptionCapability)
		preemptionVulnerability := simpolicy.PreemptionVulnerability(input.PreemptionVulnerability)
		output := simpolicy.DataNetworkConfiguration{
			AdditionalAllowedSessionTypes:       expandSimPolicyAdditionalAllowedSessionTypeResource(input.AdditionalAllowedSessionTypes),
			AllocationAndRetentionPriorityLevel: &input.AllocationAndRetentionPriorityLevel,
			DefaultSessionType:                  &defaultSessionType,
			Fiveqi:                              &input.QosIdentifier,
			PreemptionCapability:                &preemptionCapability,
			PreemptionVulnerability:             &preemptionVulnerability,
			MaximumNumberOfBufferedPackets:      &input.MaximumNumberOfBufferedPackets,
			AllowedServices:                     expandServiceResourceIdResourceModel(input.AllowedServices),
		}

		if input.DataNetworkId != "" {
			output.DataNetwork = simpolicy.DataNetworkResourceId{
				Id: input.DataNetworkId,
			}
		}

		output.SessionAmbr = expandAmbrResourceModel(input.SessionAmbr)

		outputList = append(outputList, output)
	}

	return outputList
}

func expandSimPolicyAdditionalAllowedSessionTypeResource(inputList []string) *[]simpolicy.PduSessionType {
	var outputList []simpolicy.PduSessionType
	for _, v := range inputList {
		outputList = append(outputList, simpolicy.PduSessionType(v))
	}

	return &outputList
}

func expandServiceResourceIdResourceModel(inputList []string) []simpolicy.ServiceResourceId {
	var outputList []simpolicy.ServiceResourceId
	for _, v := range inputList {
		input := v
		output := simpolicy.ServiceResourceId{
			Id: input,
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func expandAmbrResourceModel(inputList []AmbrResourceModel) simpolicy.Ambr {
	output := simpolicy.Ambr{}

	if len(inputList) == 0 {
		return output
	}

	input := inputList[0]
	output.Downlink = input.Downlink
	output.Uplink = input.Uplink

	return output
}

func flattenSliceConfigurationResourceModel(input []simpolicy.SliceConfiguration) []SliceConfigurationResourceModel {
	output := make([]SliceConfigurationResourceModel, 0)

	for _, v := range input {
		item := SliceConfigurationResourceModel{
			SliceId:              v.Slice.Id,
			DefaultDataNetworkId: v.DefaultDataNetwork.Id,
		}

		item.DataNetworkConfigurations = flattenDataNetworkConfigurationResourceModel(v.DataNetworkConfigurations)

		output = append(output, item)
	}

	return output
}

func flattenDataNetworkConfigurationResourceModel(input []simpolicy.DataNetworkConfiguration) []DataNetworkConfigurationResourceModel {
	output := make([]DataNetworkConfigurationResourceModel, 0)

	for _, v := range input {
		item := DataNetworkConfigurationResourceModel{
			DataNetworkId: v.DataNetwork.Id,
		}

		item.AdditionalAllowedSessionTypes = flattenSimPolicyAllowedSessionTypeResource(v.AdditionalAllowedSessionTypes)

		if v.AllocationAndRetentionPriorityLevel != nil {
			item.AllocationAndRetentionPriorityLevel = *v.AllocationAndRetentionPriorityLevel
		}

		item.AllowedServices = flattenServiceResourceIdResourceModel(v.AllowedServices)

		if v.DefaultSessionType != nil {
			item.DefaultSessionType = string(*v.DefaultSessionType)
		}

		if v.Fiveqi != nil {
			item.QosIdentifier = *v.Fiveqi
		}

		if v.MaximumNumberOfBufferedPackets != nil {
			item.MaximumNumberOfBufferedPackets = *v.MaximumNumberOfBufferedPackets
		}

		if v.PreemptionCapability != nil {
			item.PreemptionCapability = string(*v.PreemptionCapability)
		}

		if v.PreemptionVulnerability != nil {
			item.PreemptionVulnerability = string(*v.PreemptionVulnerability)
		}

		item.SessionAmbr = flattenSimPolicyAmbrResourceModel(v.SessionAmbr)

		output = append(output, item)
	}

	return output
}

func flattenSimPolicyAllowedSessionTypeResource(input *[]simpolicy.PduSessionType) []string {
	output := make([]string, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		output = append(output, string(v))
	}
	return output
}

func flattenServiceResourceIdResourceModel(input []simpolicy.ServiceResourceId) []string {
	output := make([]string, 0)
	for _, v := range input {
		output = append(output, v.Id)
	}
	return output
}

func flattenSimPolicyAmbrResourceModel(input simpolicy.Ambr) []AmbrResourceModel {
	return []AmbrResourceModel{
		{
			Downlink: input.Downlink,
			Uplink:   input.Uplink,
		},
	}
}
