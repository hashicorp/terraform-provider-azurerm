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

type MobileNetworkServiceDataSource struct{}

var _ sdk.DataSource = MobileNetworkServiceDataSource{}

func (r MobileNetworkServiceDataSource) ResourceType() string {
	return "azurerm_mobile_network_service"
}

func (r MobileNetworkServiceDataSource) ModelObject() interface{} {
	return &ServiceModel{}
}

func (r MobileNetworkServiceDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return service.ValidateServiceID
}

func (r MobileNetworkServiceDataSource) Arguments() map[string]*pluginsdk.Schema {
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

func (r MobileNetworkServiceDataSource) Attributes() map[string]*pluginsdk.Schema {
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

func (r MobileNetworkServiceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel ServiceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.ServiceClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(metaModel.MobileNetworkMobileNetworkId)
			if err != nil {
				return err
			}

			id := service.NewServiceID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, metaModel.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
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

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
