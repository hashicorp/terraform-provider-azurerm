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

var _ sdk.DataSource = SimPolicyDataSource{}

func (r SimPolicyDataSource) ResourceType() string {
	return "azurerm_mobile_network_sim_policy"
}

func (r SimPolicyDataSource) ModelObject() interface{} {
	return &SimPolicyModel{}
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
			var metaModel SimPolicyModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SIMPolicyClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(metaModel.MobileNetworkMobileNetworkId)
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

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
