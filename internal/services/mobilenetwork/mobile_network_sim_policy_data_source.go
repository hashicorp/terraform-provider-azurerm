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
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simpolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SimPolicyDataSource struct{}

type SimPolicyDataSourceModel struct {
	Name                string                              `tfschema:"name"`
	MobileNetworkId     string                              `tfschema:"mobile_network_id"`
	DefaultSliceId      string                              `tfschema:"default_slice_id"`
	Location            string                              `tfschema:"location"`
	RegistrationTimer   int64                               `tfschema:"registration_timer_in_seconds"`
	RfspIndex           int64                               `tfschema:"rat_frequency_selection_priority_index"`
	SliceConfigurations []SliceConfigurationDataSourceModel `tfschema:"slice"`
	Tags                map[string]string                   `tfschema:"tags"`
	UeAmbr              []AmbrDataSourceModel               `tfschema:"user_equipment_aggregate_maximum_bit_rate"`
}

type SliceConfigurationDataSourceModel struct {
	DataNetworkConfigurations []DataNetworkConfigurationDataSourceModel `tfschema:"data_network"`
	DefaultDataNetworkId      string                                    `tfschema:"default_data_network_id"`
	SliceId                   string                                    `tfschema:"slice_id"`
}

type DataNetworkConfigurationDataSourceModel struct {
	AdditionalAllowedSessionTypes       []string              `tfschema:"additional_allowed_session_types"`
	AllocationAndRetentionPriorityLevel int64                 `tfschema:"allocation_and_retention_priority_level"`
	AllowedServices                     []string              `tfschema:"allowed_services_ids"`
	DataNetworkId                       string                `tfschema:"data_network_id"`
	DefaultSessionType                  string                `tfschema:"default_session_type"`
	QosIdentifier                       int64                 `tfschema:"qos_indicator"`
	MaximumNumberOfBufferedPackets      int64                 `tfschema:"max_buffered_packets"`
	PreemptionCapability                string                `tfschema:"preemption_capability"`
	PreemptionVulnerability             string                `tfschema:"preemption_vulnerability"`
	SessionAmbr                         []AmbrDataSourceModel `tfschema:"session_aggregate_maximum_bit_rate"`
}

type AmbrDataSourceModel struct {
	Downlink string `tfschema:"downlink"`
	Uplink   string `tfschema:"uplink"`
}

var _ sdk.DataSource = SimPolicyDataSource{}

func (r SimPolicyDataSource) ResourceType() string {
	return "azurerm_mobile_network_sim_policy"
}

func (r SimPolicyDataSource) ModelObject() interface{} {
	return &SimPolicyDataSourceModel{}
}

func (r SimPolicyDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return simpolicy.ValidateSimPolicyID
}

func (r SimPolicyDataSource) Arguments() map[string]*pluginsdk.Schema {
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

func (r SimPolicyDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"default_slice_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"registration_timer_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"rat_frequency_selection_priority_index": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"slice": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_network": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"additional_allowed_session_types": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"allocation_and_retention_priority_level": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"allowed_services_ids": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"data_network_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"default_session_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"qos_indicator": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
								"max_buffered_packets": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"preemption_capability": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"preemption_vulnerability": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"session_aggregate_maximum_bit_rate": {
									Type:     pluginsdk.TypeList,
									Computed: true,

									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"downlink": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"uplink": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"default_data_network_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"slice_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"user_equipment_aggregate_maximum_bit_rate": {
			Type:     pluginsdk.TypeList,
			Computed: true,

			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"downlink": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"uplink": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r SimPolicyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel SimPolicyDataSourceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SIMPolicyClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(metaModel.MobileNetworkId)
			if err != nil {
				return err
			}

			id := simpolicy.NewSimPolicyID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, metaModel.Name)

			resp, err := client.SimPoliciesGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := SimPolicyDataSourceModel{
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

				state.SliceConfigurations = flattenSliceConfigurationDataSourceModel(model.Properties.SliceConfigurations)
				state.UeAmbr = flattenSimPolicyAmbrDataSourceModel(model.Properties.UeAmbr)

				if model.Tags != nil {
					state.Tags = *model.Tags
				}

			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func flattenSliceConfigurationDataSourceModel(input []simpolicy.SliceConfiguration) []SliceConfigurationDataSourceModel {
	output := make([]SliceConfigurationDataSourceModel, 0)

	for _, v := range input {
		output = append(output, SliceConfigurationDataSourceModel{
			DataNetworkConfigurations: flattenDataNetworkConfigurationDataSourceModel(v.DataNetworkConfigurations),
			DefaultDataNetworkId:      v.DefaultDataNetwork.Id,
			SliceId:                   v.Slice.Id,
		})
	}

	return output
}

func flattenDataNetworkConfigurationDataSourceModel(inputList []simpolicy.DataNetworkConfiguration) []DataNetworkConfigurationDataSourceModel {
	output := make([]DataNetworkConfigurationDataSourceModel, 0)

	for _, input := range inputList {
		item := DataNetworkConfigurationDataSourceModel{
			DataNetworkId: input.DataNetwork.Id,
		}

		item.AdditionalAllowedSessionTypes = flattenSimPolicyAllowedSessionTypeDataSource(input.AdditionalAllowedSessionTypes)

		if input.AllocationAndRetentionPriorityLevel != nil {
			item.AllocationAndRetentionPriorityLevel = *input.AllocationAndRetentionPriorityLevel
		}

		item.AllowedServices = flattenServiceResourceIdDataSourceModel(input.AllowedServices)

		if input.DefaultSessionType != nil {
			item.DefaultSessionType = string(*input.DefaultSessionType)
		}

		if input.Fiveqi != nil {
			item.QosIdentifier = *input.Fiveqi
		}

		if input.MaximumNumberOfBufferedPackets != nil {
			item.MaximumNumberOfBufferedPackets = *input.MaximumNumberOfBufferedPackets
		}

		if input.PreemptionCapability != nil {
			item.PreemptionCapability = string(*input.PreemptionCapability)
		}

		if input.PreemptionVulnerability != nil {
			item.PreemptionVulnerability = string(*input.PreemptionVulnerability)
		}

		item.SessionAmbr = flattenSimPolicyAmbrDataSourceModel(input.SessionAmbr)

		output = append(output, item)
	}

	return output
}

func flattenSimPolicyAllowedSessionTypeDataSource(input *[]simpolicy.PduSessionType) []string {
	output := make([]string, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		output = append(output, string(v))
	}
	return output
}

func flattenServiceResourceIdDataSourceModel(input []simpolicy.ServiceResourceId) []string {
	output := make([]string, 0)

	for _, v := range input {
		output = append(output, v.Id)
	}

	return output
}

func flattenSimPolicyAmbrDataSourceModel(input simpolicy.Ambr) []AmbrDataSourceModel {
	return []AmbrDataSourceModel{
		{
			Downlink: input.Downlink,
			Uplink:   input.Uplink,
		},
	}
}
