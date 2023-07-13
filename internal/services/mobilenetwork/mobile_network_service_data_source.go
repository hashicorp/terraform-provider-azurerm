// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ServiceDataSource struct{}

type ServiceDataSourceModel struct {
	Name              string                                       `tfschema:"name"`
	MobileNetworkId   string                                       `tfschema:"mobile_network_id"`
	Location          string                                       `tfschema:"location"`
	PccRules          []ServiceDataSourcePccRuleConfigurationModel `tfschema:"pcc_rule"`
	ServicePrecedence int64                                        `tfschema:"service_precedence"`
	ServiceQosPolicy  []ServiceDataSourceQosPolicyModel            `tfschema:"service_qos_policy"`
	Tags              map[string]string                            `tfschema:"tags"`
}

type ServiceDataSourcePccRuleConfigurationModel struct {
	RuleName                 string                                          `tfschema:"name"`
	RulePrecedence           int64                                           `tfschema:"precedence"`
	RuleQosPolicy            []ServiceDataSourcePccRuleQosPolicyModel        `tfschema:"qos_policy"`
	ServiceDataFlowTemplates []ServiceDataSourceServiceDataFlowTemplateModel `tfschema:"service_data_flow_template"`
	TrafficControlEnabled    bool                                            `tfschema:"traffic_control_enabled"`
}

type ServiceDataSourcePccRuleQosPolicyModel struct {
	AllocationAndRetentionPriorityLevel int64                           `tfschema:"allocation_and_retention_priority_level"`
	QosIdentifier                       int64                           `tfschema:"qos_indicator"`
	GuaranteedBitRate                   []ServiceDataSourceBitRateModel `tfschema:"guaranteed_bit_rate"`
	MaximumBitRate                      []ServiceDataSourceBitRateModel `tfschema:"maximum_bit_rate"`
	PreemptionCapability                string                          `tfschema:"preemption_capability"`
	PreemptionVulnerability             string                          `tfschema:"preemption_vulnerability"`
}

type ServiceDataSourceBitRateModel struct {
	Downlink string `tfschema:"downlink"`
	Uplink   string `tfschema:"uplink"`
}

type ServiceDataSourceServiceDataFlowTemplateModel struct {
	Direction    string   `tfschema:"direction"`
	Ports        []string `tfschema:"ports"`
	Protocol     []string `tfschema:"protocol"`
	RemoteIPList []string `tfschema:"remote_ip_list"`
	TemplateName string   `tfschema:"name"`
}

type ServiceDataSourceQosPolicyModel struct {
	AllocationAndRetentionPriorityLevel int64                           `tfschema:"allocation_and_retention_priority_level"`
	QosIdentifier                       int64                           `tfschema:"qos_indicator"`
	MaximumBitRate                      []ServiceDataSourceBitRateModel `tfschema:"maximum_bit_rate"`
	PreemptionCapability                string                          `tfschema:"preemption_capability"`
	PreemptionVulnerability             string                          `tfschema:"preemption_vulnerability"`
}

var _ sdk.DataSource = ServiceDataSource{}

func (r ServiceDataSource) ResourceType() string {
	return "azurerm_mobile_network_service"
}

func (r ServiceDataSource) ModelObject() interface{} {
	return &ServiceDataSourceModel{}
}

func (r ServiceDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return service.ValidateServiceID
}

func (r ServiceDataSource) Arguments() map[string]*pluginsdk.Schema {
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

func (r ServiceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"location": commonschema.LocationComputed(),

		"pcc_rule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"precedence": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"qos_policy": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allocation_and_retention_priority_level": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"qos_indicator": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"guaranteed_bit_rate": {
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

								"maximum_bit_rate": {
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

								"preemption_capability": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"preemption_vulnerability": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"service_data_flow_template": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"direction": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"ports": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"protocol": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"remote_ip_list": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"traffic_control_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},

		"service_precedence": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"service_qos_policy": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_and_retention_priority_level": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"qos_indicator": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"maximum_bit_rate": {
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

					"preemption_capability": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"preemption_vulnerability": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r ServiceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel ServiceDataSourceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.ServiceClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(metaModel.MobileNetworkId)
			if err != nil {
				return err
			}

			id := service.NewServiceID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, metaModel.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)
			state := ServiceDataSourceModel{
				Name:            id.ServiceName,
				MobileNetworkId: mobilenetwork.NewMobileNetworkID(id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName).ID(),
			}
			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				props := model.Properties
				state.PccRules = flattenPccRuleConfigurationDataSourceModel(props.PccRules)
				state.ServicePrecedence = props.ServicePrecedence
				state.ServiceQosPolicy = flattenQosPolicyDataSourceModel(props.ServiceQosPolicy)

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenPccRuleConfigurationDataSourceModel(inputList []service.PccRuleConfiguration) []ServiceDataSourcePccRuleConfigurationModel {
	var outputList []ServiceDataSourcePccRuleConfigurationModel

	for _, input := range inputList {
		output := ServiceDataSourcePccRuleConfigurationModel{
			RuleName:       input.RuleName,
			RulePrecedence: input.RulePrecedence,
		}

		output.RuleQosPolicy = flattenPccRuleQosPolicyDataSourceModel(input.RuleQosPolicy)

		output.ServiceDataFlowTemplates = flattenServiceDataFlowTemplateDataSourceModel(&input.ServiceDataFlowTemplates)

		if input.TrafficControl != nil {
			output.TrafficControlEnabled = *input.TrafficControl == service.TrafficControlPermissionEnabled
		}

		outputList = append(outputList, output)
	}

	return outputList
}
func flattenPccRuleQosPolicyDataSourceModel(input *service.PccRuleQosPolicy) []ServiceDataSourcePccRuleQosPolicyModel {
	if input == nil {
		return []ServiceDataSourcePccRuleQosPolicyModel{}
	}

	output := ServiceDataSourcePccRuleQosPolicyModel{}

	if input.AllocationAndRetentionPriorityLevel != nil {
		output.AllocationAndRetentionPriorityLevel = *input.AllocationAndRetentionPriorityLevel
	}

	if input.Fiveqi != nil {
		output.QosIdentifier = *input.Fiveqi
	}

	output.GuaranteedBitRate = flattenBitRateDataSourceModel(input.GuaranteedBitRate)

	output.MaximumBitRate = flattenBitRateDataSourceModel(&input.MaximumBitRate)

	if input.PreemptionCapability != nil {
		output.PreemptionCapability = string(*input.PreemptionCapability)
	}

	if input.PreemptionVulnerability != nil {
		output.PreemptionVulnerability = string(*input.PreemptionVulnerability)
	}

	return []ServiceDataSourcePccRuleQosPolicyModel{
		output,
	}
}

func flattenServiceDataFlowTemplateDataSourceModel(inputList *[]service.ServiceDataFlowTemplate) []ServiceDataSourceServiceDataFlowTemplateModel {
	output := make([]ServiceDataSourceServiceDataFlowTemplateModel, 0)
	if inputList == nil {
		return output
	}

	for _, input := range *inputList {
		model := ServiceDataSourceServiceDataFlowTemplateModel{
			Direction:    string(input.Direction),
			Protocol:     input.Protocol,
			RemoteIPList: input.RemoteIPList,
			TemplateName: input.TemplateName,
		}

		if input.Ports != nil {
			model.Ports = *input.Ports
		}
		output = append(output, model)
	}

	return output
}

func flattenBitRateDataSourceModel(input *service.Ambr) []ServiceDataSourceBitRateModel {
	if input == nil {
		return []ServiceDataSourceBitRateModel{}
	}

	return []ServiceDataSourceBitRateModel{
		{
			Downlink: input.Downlink,
			Uplink:   input.Uplink,
		},
	}
}

func flattenQosPolicyDataSourceModel(input *service.QosPolicy) []ServiceDataSourceQosPolicyModel {
	if input == nil {
		return []ServiceDataSourceQosPolicyModel{}
	}

	output := ServiceDataSourceQosPolicyModel{}

	if input.AllocationAndRetentionPriorityLevel != nil {
		output.AllocationAndRetentionPriorityLevel = *input.AllocationAndRetentionPriorityLevel
	}

	if input.Fiveqi != nil {
		output.QosIdentifier = *input.Fiveqi
	}

	output.MaximumBitRate = flattenBitRateDataSourceModel(&input.MaximumBitRate)

	if input.PreemptionCapability != nil {
		output.PreemptionCapability = string(*input.PreemptionCapability)
	}

	if input.PreemptionVulnerability != nil {
		output.PreemptionVulnerability = string(*input.PreemptionVulnerability)
	}

	return []ServiceDataSourceQosPolicyModel{
		output,
	}
}
