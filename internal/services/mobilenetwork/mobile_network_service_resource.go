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
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ServiceModel struct {
	Name                         string                      `tfschema:"name"`
	MobileNetworkMobileNetworkId string                      `tfschema:"mobile_network_id"`
	Location                     string                      `tfschema:"location"`
	PccRules                     []PccRuleConfigurationModel `tfschema:"pcc_rule"`
	ServicePrecedence            int64                       `tfschema:"service_precedence"`
	ServiceQosPolicy             []QosPolicyModel            `tfschema:"service_qos_policy"`
	Tags                         map[string]string           `tfschema:"tags"`
}

type PccRuleConfigurationModel struct {
	RuleName                 string                         `tfschema:"name"`
	RulePrecedence           int64                          `tfschema:"precedence"`
	RuleQosPolicy            []PccRuleQosPolicyModel        `tfschema:"qos_policy"`
	ServiceDataFlowTemplates []ServiceDataFlowTemplateModel `tfschema:"service_data_flow_template"`
	TrafficControlEnabled    bool                           `tfschema:"traffic_control_enabled"`
}

type PccRuleQosPolicyModel struct {
	AllocationAndRetentionPriorityLevel int64          `tfschema:"allocation_and_retention_priority_level"`
	QosIdentifier                       int64          `tfschema:"qos_indicator"`
	GuaranteedBitRate                   []BitRateModel `tfschema:"guaranteed_bit_rate"`
	MaximumBitRate                      []BitRateModel `tfschema:"maximum_bit_rate"`
	PreemptionCapability                string         `tfschema:"preemption_capability"`
	PreemptionVulnerability             string         `tfschema:"preemption_vulnerability"`
}

type BitRateModel struct {
	Downlink string `tfschema:"downlink"`
	Uplink   string `tfschema:"uplink"`
}

type ServiceDataFlowTemplateModel struct {
	Direction    string   `tfschema:"direction"`
	Ports        []string `tfschema:"ports"`
	Protocol     []string `tfschema:"protocol"`
	RemoteIPList []string `tfschema:"remote_ip_list"`
	TemplateName string   `tfschema:"name"`
}

type QosPolicyModel struct {
	AllocationAndRetentionPriorityLevel int64          `tfschema:"allocation_and_retention_priority_level"`
	QosIdentifier                       int64          `tfschema:"qos_indicator"`
	MaximumBitRate                      []BitRateModel `tfschema:"maximum_bit_rate"`
	PreemptionCapability                string         `tfschema:"preemption_capability"`
	PreemptionVulnerability             string         `tfschema:"preemption_vulnerability"`
}

type MobileNetworkServiceResource struct{}

var _ sdk.ResourceWithUpdate = MobileNetworkServiceResource{}

func (r MobileNetworkServiceResource) ResourceType() string {
	return "azurerm_mobile_network_service"
}

func (r MobileNetworkServiceResource) ModelObject() interface{} {
	return &ServiceModel{}
}

func (r MobileNetworkServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return service.ValidateServiceID
}

func (r MobileNetworkServiceResource) Arguments() map[string]*pluginsdk.Schema {
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

		"location": commonschema.Location(),

		"pcc_rule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 64),
					},

					"precedence": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 255),
					},

					"qos_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allocation_and_retention_priority_level": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},

								"qos_indicator": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},

								"guaranteed_bit_rate": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"downlink": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringMatch(
													regexp.MustCompile(`^([1-9]\d*|0)(\.\d+)?\s(Kbps|Mbps|Gbps|Tbps)$`),
													"The value must be a number followed by Kbps, Mbps, Gbps or Tbps.",
												)},

											"uplink": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringMatch(
													regexp.MustCompile(`^([1-9]\d*|0)(\.\d+)?\s(Kbps|Mbps|Gbps|Tbps)$`),
													"The value must be a number followed by Kbps, Mbps, Gbps or Tbps.",
												)},
										},
									},
								},

								"maximum_bit_rate": {
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
												)},

											"uplink": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringMatch(
													regexp.MustCompile(`^([1-9]\d*|0)(\.\d+)?\s(Kbps|Mbps|Gbps|Tbps)$`),
													"The value must be a number followed by Kbps, Mbps, Gbps or Tbps.",
												)},
										},
									},
								},

								"preemption_capability": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(service.PreemptionCapabilityNotPreempt),
									ValidateFunc: validation.StringInSlice([]string{
										string(service.PreemptionCapabilityNotPreempt),
										string(service.PreemptionCapabilityMayPreempt),
									}, false),
								},

								"preemption_vulnerability": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(service.PreemptionVulnerabilityPreemptable),
									ValidateFunc: validation.StringInSlice([]string{
										string(service.PreemptionVulnerabilityNotPreemptable),
										string(service.PreemptionVulnerabilityPreemptable),
									}, false),
								},
							},
						},
					},

					"service_data_flow_template": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"direction": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(service.SdfDirectionUplink),
										string(service.SdfDirectionDownlink),
										string(service.SdfDirectionBidirectional),
									}, false),
								},

								"ports": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"protocol": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"remote_ip_list": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"traffic_control_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"service_precedence": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 255),
		},

		"service_qos_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allocation_and_retention_priority_level": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      9,
						ValidateFunc: validation.IntBetween(1, 127),
					},

					"qos_indicator": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"maximum_bit_rate": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"downlink": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^\d+(\\.\\d+)?\s(bps|Kbps|Mbps|Gbps|Tbps)$`),
										"The value must be a number followed by bps, Kbps, Mbps, Gbps or Tbps.",
									)},

								"uplink": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^\d+(\\.\\d+)?\s(bps|Kbps|Mbps|Gbps|Tbps)$`),
										"The value must be a number followed by bps, Kbps, Mbps, Gbps or Tbps.",
									)},
							},
						},
					},

					"preemption_capability": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(service.PreemptionCapabilityNotPreempt),
							string(service.PreemptionCapabilityMayPreempt),
						}, false),
					},

					"preemption_vulnerability": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(service.PreemptionVulnerabilityNotPreemptable),
							string(service.PreemptionVulnerabilityPreemptable),
						}, false),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r MobileNetworkServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MobileNetworkServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.ServiceClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(model.MobileNetworkMobileNetworkId)
			if err != nil {
				return err
			}

			id := service.NewServiceID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := service.Service{
				Location: location.Normalize(model.Location),
				Properties: service.ServicePropertiesFormat{
					ServicePrecedence: model.ServicePrecedence,
					PccRules:          expandPccRuleConfigurationModel(model.PccRules),
					ServiceQosPolicy:  expandQosPolicyModel(model.ServiceQosPolicy),
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MobileNetworkServiceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.ServiceClient

			id, err := service.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ServiceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			properties := *resp.Model

			if metadata.ResourceData.HasChange("pcc_rule") {
				properties.Properties.PccRules = expandPccRuleConfigurationModel(model.PccRules)
			}

			if metadata.ResourceData.HasChange("service_precedence") {
				properties.Properties.ServicePrecedence = model.ServicePrecedence
			}

			if metadata.ResourceData.HasChange("service_qos_policy") {
				properties.Properties.ServiceQosPolicy = expandQosPolicyModel(model.ServiceQosPolicy)
			}

			properties.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MobileNetworkServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.ServiceClient

			id, err := service.ParseServiceID(metadata.ResourceData.Id())
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

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			model := *resp.Model

			state := ServiceModel{
				Name:                         id.ServiceName,
				MobileNetworkMobileNetworkId: mobilenetwork.NewMobileNetworkID(id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName).ID(),
				Location:                     location.Normalize(model.Location),
			}

			properties := model.Properties

			state.PccRules = flattenPccRuleConfigurationModel(properties.PccRules)

			state.ServicePrecedence = properties.ServicePrecedence

			state.ServiceQosPolicy = flattenQosPolicyModel(properties.ServiceQosPolicy)

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MobileNetworkServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.ServiceClient

			id, err := service.ParseServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := resourceMobileNetworkChildWaitForDeletion(ctx, id.ID(), func() (*http.Response, error) {
				resp, err := client.Get(ctx, *id)
				return resp.HttpResponse, err
			}); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandPccRuleConfigurationModel(inputList []PccRuleConfigurationModel) []service.PccRuleConfiguration {
	var outputList []service.PccRuleConfiguration
	for _, v := range inputList {
		input := v
		output := service.PccRuleConfiguration{
			RuleName:       input.RuleName,
			RulePrecedence: input.RulePrecedence,
		}

		trafficControlValue := service.TrafficControlPermissionBlocked
		if input.TrafficControlEnabled {
			trafficControlValue = service.TrafficControlPermissionEnabled
		}
		output.TrafficControl = &trafficControlValue

		output.RuleQosPolicy = expandPccRuleQosPolicyModel(input.RuleQosPolicy)

		output.ServiceDataFlowTemplates = expandServiceDataFlowTemplateModel(input.ServiceDataFlowTemplates)

		outputList = append(outputList, output)
	}

	return outputList
}

func expandPccRuleQosPolicyModel(inputList []PccRuleQosPolicyModel) *service.PccRuleQosPolicy {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	capability := service.PreemptionCapability(input.PreemptionCapability)
	vulnerability := service.PreemptionVulnerability(input.PreemptionVulnerability)
	output := service.PccRuleQosPolicy{
		AllocationAndRetentionPriorityLevel: &input.AllocationAndRetentionPriorityLevel,
		Fiveqi:                              &input.QosIdentifier,
		PreemptionCapability:                &capability,
		PreemptionVulnerability:             &vulnerability,
	}

	output.GuaranteedBitRate = expandBitRateModel(input.GuaranteedBitRate)

	if v := expandBitRateModel(input.MaximumBitRate); v != nil {
		output.MaximumBitRate = *v
	}

	return &output
}

func expandServiceDataFlowTemplateModel(inputList []ServiceDataFlowTemplateModel) []service.ServiceDataFlowTemplate {
	outputList := make([]service.ServiceDataFlowTemplate, 0)
	for _, v := range inputList {
		input := v
		output := service.ServiceDataFlowTemplate{
			Direction:    service.SdfDirection(input.Direction),
			Ports:        &input.Ports,
			Protocol:     input.Protocol,
			RemoteIPList: input.RemoteIPList,
			TemplateName: input.TemplateName,
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func expandQosPolicyModel(inputList []QosPolicyModel) *service.QosPolicy {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	capability := service.PreemptionCapability(input.PreemptionCapability)
	vulnerability := service.PreemptionVulnerability(input.PreemptionVulnerability)
	output := service.QosPolicy{
		AllocationAndRetentionPriorityLevel: &input.AllocationAndRetentionPriorityLevel,
		Fiveqi:                              &input.QosIdentifier,
		PreemptionCapability:                &capability,
		PreemptionVulnerability:             &vulnerability,
	}

	if v := expandBitRateModel(input.MaximumBitRate); v != nil {
		output.MaximumBitRate = *v
	}

	return &output
}

func flattenPccRuleConfigurationModel(inputList []service.PccRuleConfiguration) []PccRuleConfigurationModel {
	var outputList []PccRuleConfigurationModel

	for _, input := range inputList {
		output := PccRuleConfigurationModel{
			RuleName:       input.RuleName,
			RulePrecedence: input.RulePrecedence,
		}

		output.RuleQosPolicy = flattenPccRuleQosPolicyModel(input.RuleQosPolicy)

		output.ServiceDataFlowTemplates = flattenServiceDataFlowTemplateModel(&input.ServiceDataFlowTemplates)

		if input.TrafficControl != nil {
			output.TrafficControlEnabled = *input.TrafficControl == service.TrafficControlPermissionEnabled
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenPccRuleQosPolicyModel(input *service.PccRuleQosPolicy) []PccRuleQosPolicyModel {
	var outputList []PccRuleQosPolicyModel
	if input == nil {
		return outputList
	}

	output := PccRuleQosPolicyModel{}

	if input.AllocationAndRetentionPriorityLevel != nil {
		output.AllocationAndRetentionPriorityLevel = *input.AllocationAndRetentionPriorityLevel
	}

	if input.Fiveqi != nil {
		output.QosIdentifier = *input.Fiveqi
	}

	output.GuaranteedBitRate = flattenBitRateModel(input.GuaranteedBitRate)

	output.MaximumBitRate = flattenBitRateModel(&input.MaximumBitRate)

	if input.PreemptionCapability != nil {
		output.PreemptionCapability = string(*input.PreemptionCapability)
	}

	if input.PreemptionVulnerability != nil {
		output.PreemptionVulnerability = string(*input.PreemptionVulnerability)
	}

	return append(outputList, output)
}

func flattenServiceDataFlowTemplateModel(inputList *[]service.ServiceDataFlowTemplate) []ServiceDataFlowTemplateModel {
	var outputList []ServiceDataFlowTemplateModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := ServiceDataFlowTemplateModel{
			Direction:    string(input.Direction),
			Protocol:     input.Protocol,
			RemoteIPList: input.RemoteIPList,
			TemplateName: input.TemplateName,
		}

		if input.Ports != nil {
			output.Ports = *input.Ports
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenQosPolicyModel(input *service.QosPolicy) []QosPolicyModel {
	var outputList []QosPolicyModel
	if input == nil {
		return outputList
	}

	output := QosPolicyModel{}

	if input.AllocationAndRetentionPriorityLevel != nil {
		output.AllocationAndRetentionPriorityLevel = *input.AllocationAndRetentionPriorityLevel
	}

	if input.Fiveqi != nil {
		output.QosIdentifier = *input.Fiveqi
	}

	output.MaximumBitRate = flattenBitRateModel(&input.MaximumBitRate)

	if input.PreemptionCapability != nil {
		output.PreemptionCapability = string(*input.PreemptionCapability)
	}

	if input.PreemptionVulnerability != nil {
		output.PreemptionVulnerability = string(*input.PreemptionVulnerability)
	}

	return append(outputList, output)
}

// make it return a pointer because some property accept nil value
func expandBitRateModel(inputList []BitRateModel) *service.Ambr {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := service.Ambr{
		Downlink: input.Downlink,
		Uplink:   input.Uplink,
	}

	return &output
}

func flattenBitRateModel(input *service.Ambr) []BitRateModel {
	var outputList []BitRateModel
	if input == nil {
		return outputList
	}

	output := BitRateModel{
		Downlink: input.Downlink,
		Uplink:   input.Uplink,
	}

	return append(outputList, output)
}
