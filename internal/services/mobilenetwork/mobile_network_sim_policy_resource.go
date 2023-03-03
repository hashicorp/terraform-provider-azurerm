package mobilenetwork

import (
	"context"
	"fmt"
	"net/http"
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

type SimPolicyModel struct {
	Name                         string                    `tfschema:"name"`
	MobileNetworkMobileNetworkId string                    `tfschema:"mobile_network_id"`
	DefaultSliceId               string                    `tfschema:"default_slice_id"`
	Location                     string                    `tfschema:"location"`
	RegistrationTimer            int64                     `tfschema:"registration_timer_in_seconds"`
	RfspIndex                    int64                     `tfschema:"rat_frequency_selection_priority_index"`
	SliceConfigurations          []SliceConfigurationModel `tfschema:"slice"`
	Tags                         map[string]string         `tfschema:"tags"`
	UeAmbr                       []AmbrModel               `tfschema:"user_equipment_aggregate_maximum_bit_rate"`
}

type SliceConfigurationModel struct {
	DataNetworkConfigurations []DataNetworkConfigurationModel `tfschema:"data_network"`
	DefaultDataNetworkId      string                          `tfschema:"default_data_network_id"`
	SliceId                   string                          `tfschema:"slice_id"`
}

type DataNetworkConfigurationModel struct {
	AdditionalAllowedSessionTypes       []string    `tfschema:"additional_allowed_session_types"`
	AllocationAndRetentionPriorityLevel int64       `tfschema:"allocation_and_retention_priority_level"`
	AllowedServices                     []string    `tfschema:"allowed_services_ids"`
	DataNetworkId                       string      `tfschema:"data_network_id"`
	DefaultSessionType                  string      `tfschema:"default_session_type"`
	QosIdentifier                       int64       `tfschema:"qos_indicator"`
	MaximumNumberOfBufferedPackets      int64       `tfschema:"max_buffered_packets"`
	PreemptionCapability                string      `tfschema:"preemption_capability"`
	PreemptionVulnerability             string      `tfschema:"preemption_vulnerability"`
	SessionAmbr                         []AmbrModel `tfschema:"session_aggregate_maximum_bit_rate"`
}

type AmbrModel struct {
	Downlink string `tfschema:"downlink"`
	Uplink   string `tfschema:"uplink"`
}

type SimPolicyResource struct{}

var _ sdk.ResourceWithUpdate = SimPolicyResource{}

func (r SimPolicyResource) ResourceType() string {
	return "azurerm_mobile_network_sim_policy"
}

func (r SimPolicyResource) ModelObject() interface{} {
	return &SimPolicyModel{}
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
			var model SimPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SIMPolicyClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(model.MobileNetworkMobileNetworkId)
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

			properties.Properties.SliceConfigurations = expandSliceConfigurationModel(model.SliceConfigurations)

			properties.Properties.UeAmbr = expandAmbrModel(model.UeAmbr)

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

			var plan SimPolicyModel
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
				model.Properties.SliceConfigurations = expandSliceConfigurationModel(plan.SliceConfigurations)
			}

			if metadata.ResourceData.HasChange("user_equipment_aggregate_maximum_bit_rate") {
				model.Properties.UeAmbr = expandAmbrModel(plan.UeAmbr)
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

			state := SimPolicyModel{
				Name:                         id.SimPolicyName,
				MobileNetworkMobileNetworkId: mobilenetwork.NewMobileNetworkID(id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName).ID(),
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

				state.SliceConfigurations = flattenSliceConfigurationModel(model.Properties.SliceConfigurations)

				state.UeAmbr = flattenSimPolicyAmbrModel(model.Properties.UeAmbr)

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

			if err := resourceMobileNetworkChildWaitForDeletion(ctx, id.ID(), func() (*http.Response, error) {
				resp, err := client.SimPoliciesGet(ctx, *id)
				return resp.HttpResponse, err
			}); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandSliceConfigurationModel(inputList []SliceConfigurationModel) []simpolicy.SliceConfiguration {
	var outputList []simpolicy.SliceConfiguration
	for _, v := range inputList {
		input := v
		output := simpolicy.SliceConfiguration{
			Slice: simpolicy.SliceResourceId{Id: input.SliceId},
		}

		output.DataNetworkConfigurations = expandDataNetworkConfigurationModel(input.DataNetworkConfigurations)

		if input.DefaultDataNetworkId != "" {
			output.DefaultDataNetwork = simpolicy.DataNetworkResourceId{
				Id: input.DefaultDataNetworkId,
			}
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func expandDataNetworkConfigurationModel(inputList []DataNetworkConfigurationModel) []simpolicy.DataNetworkConfiguration {
	var outputList []simpolicy.DataNetworkConfiguration
	for _, v := range inputList {
		input := v

		defaultSessionType := simpolicy.PduSessionType(input.DefaultSessionType)
		preemptionCapability := simpolicy.PreemptionCapability(input.PreemptionCapability)
		preemptionVulnerability := simpolicy.PreemptionVulnerability(input.PreemptionVulnerability)
		output := simpolicy.DataNetworkConfiguration{
			AdditionalAllowedSessionTypes:       expandSimPolicyAdditionalAllowedSessionType(input.AdditionalAllowedSessionTypes),
			AllocationAndRetentionPriorityLevel: &input.AllocationAndRetentionPriorityLevel,
			DefaultSessionType:                  &defaultSessionType,
			Fiveqi:                              &input.QosIdentifier,
			PreemptionCapability:                &preemptionCapability,
			PreemptionVulnerability:             &preemptionVulnerability,
			MaximumNumberOfBufferedPackets:      &input.MaximumNumberOfBufferedPackets,
			AllowedServices:                     expandServiceResourceIdModel(input.AllowedServices),
		}

		if input.DataNetworkId != "" {
			output.DataNetwork = simpolicy.DataNetworkResourceId{
				Id: input.DataNetworkId,
			}
		}

		output.SessionAmbr = expandAmbrModel(input.SessionAmbr)

		outputList = append(outputList, output)
	}

	return outputList
}

func expandSimPolicyAdditionalAllowedSessionType(inputList []string) *[]simpolicy.PduSessionType {
	var outputList []simpolicy.PduSessionType
	for _, v := range inputList {
		outputList = append(outputList, simpolicy.PduSessionType(v))
	}

	return &outputList
}

func expandServiceResourceIdModel(inputList []string) []simpolicy.ServiceResourceId {
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

func expandAmbrModel(inputList []AmbrModel) simpolicy.Ambr {
	output := simpolicy.Ambr{}

	if len(inputList) == 0 {
		return output
	}

	input := &inputList[0]
	output.Downlink = input.Downlink
	output.Uplink = input.Uplink

	return output
}

func flattenSliceConfigurationModel(inputList []simpolicy.SliceConfiguration) []SliceConfigurationModel {
	var outputList []SliceConfigurationModel

	for _, input := range inputList {
		output := SliceConfigurationModel{
			SliceId:              input.Slice.Id,
			DefaultDataNetworkId: input.DefaultDataNetwork.Id,
		}

		output.DataNetworkConfigurations = flattenDataNetworkConfigurationModel(&input.DataNetworkConfigurations)

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenDataNetworkConfigurationModel(inputList *[]simpolicy.DataNetworkConfiguration) []DataNetworkConfigurationModel {
	var outputList []DataNetworkConfigurationModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := DataNetworkConfigurationModel{
			DataNetworkId: input.DataNetwork.Id,
		}

		output.AdditionalAllowedSessionTypes = flattenSimPolicyAllowedSessionType(input.AdditionalAllowedSessionTypes)

		if input.AllocationAndRetentionPriorityLevel != nil {
			output.AllocationAndRetentionPriorityLevel = *input.AllocationAndRetentionPriorityLevel
		}

		output.AllowedServices = flattenServiceResourceIdModel(&input.AllowedServices)

		if input.DefaultSessionType != nil {
			output.DefaultSessionType = string(*input.DefaultSessionType)
		}

		if input.Fiveqi != nil {
			output.QosIdentifier = *input.Fiveqi
		}

		if input.MaximumNumberOfBufferedPackets != nil {
			output.MaximumNumberOfBufferedPackets = *input.MaximumNumberOfBufferedPackets
		}

		if input.PreemptionCapability != nil {
			output.PreemptionCapability = string(*input.PreemptionCapability)
		}

		if input.PreemptionVulnerability != nil {
			output.PreemptionVulnerability = string(*input.PreemptionVulnerability)
		}

		output.SessionAmbr = flattenSimPolicyAmbrModel(input.SessionAmbr)

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenSimPolicyAllowedSessionType(input *[]simpolicy.PduSessionType) []string {
	output := make([]string, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		output = append(output, string(v))
	}
	return output
}

func flattenServiceResourceIdModel(inputList *[]simpolicy.ServiceResourceId) []string {
	var outputList []string
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := input.Id
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenSimPolicyAmbrModel(input simpolicy.Ambr) []AmbrModel {
	var outputList []AmbrModel

	output := AmbrModel{
		Downlink: input.Downlink,
		Uplink:   input.Uplink,
	}

	return append(outputList, output)
}
