package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/sim"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SimDataSourceModel struct {
	Name                                  string                                 `tfschema:"name"`
	MobileNetworkSimGroupId               string                                 `tfschema:"mobile_network_sim_group_id"`
	DeviceType                            string                                 `tfschema:"device_type"`
	IntegratedCircuitCardIdentifier       string                                 `tfschema:"integrated_circuit_card_identifier"`
	InternationalMobileSubscriberIdentity string                                 `tfschema:"international_mobile_subscriber_identity"`
	SimPolicyId                           string                                 `tfschema:"sim_policy_id"`
	StaticIPConfiguration                 []SimStaticIPPropertiesDataSourceModel `tfschema:"static_ip_configuration"`
	SimState                              string                                 `tfschema:"sim_state"`
	VendorKeyFingerprint                  string                                 `tfschema:"vendor_key_fingerprint"`
	VendorName                            string                                 `tfschema:"vendor_name"`
}

type SimStaticIPPropertiesDataSourceModel struct {
	AttachedDataNetworkId string `tfschema:"attached_data_network_id"`
	SliceId               string `tfschema:"slice_id"`
	StaticIP              string `tfschema:"static_ipv4_address"`
}

type SimDataSource struct{}

var _ sdk.DataSource = SimDataSource{}

func (r SimDataSource) ResourceType() string {
	return "azurerm_mobile_network_sim"
}

func (r SimDataSource) ModelObject() interface{} {
	return &SimDataSourceModel{}
}

func (r SimDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sim.ValidateSimID
}

func (r SimDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_sim_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: simgroup.ValidateSimGroupID,
		},
	}
}

func (r SimDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"integrated_circuit_card_identifier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"international_mobile_subscriber_identity": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"device_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sim_policy_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"static_ip_configuration": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"attached_data_network_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"slice_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"static_ipv4_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"sim_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vendor_key_fingerprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vendor_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r SimDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel SimResourceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SIMClient
			simGroupId, err := simgroup.ParseSimGroupID(metaModel.MobileNetworkSimGroupId)
			if err != nil {
				return err
			}

			id := sim.NewSimID(simGroupId.SubscriptionId, simGroupId.ResourceGroupName, simGroupId.SimGroupName, metaModel.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := SimDataSourceModel{
				Name:                    id.SimName,
				MobileNetworkSimGroupId: simgroup.NewSimGroupID(id.SubscriptionId, id.ResourceGroupName, id.SimGroupName).ID(),
			}

			if model := resp.Model; model != nil {
				prop := model.Properties

				state.DeviceType = pointer.From(prop.DeviceType)
				state.IntegratedCircuitCardIdentifier = pointer.From(prop.IntegratedCircuitCardIdentifier)
				state.SimState = string(pointer.From(prop.SimState))
				state.InternationalMobileSubscriberIdentity = prop.InternationalMobileSubscriberIdentity
				state.VendorKeyFingerprint = pointer.From(prop.VendorKeyFingerprint)
				state.VendorName = pointer.From(prop.VendorName)
				if simPolicy := prop.SimPolicy; prop.SimPolicy != nil {
					state.SimPolicyId = simPolicy.Id
				}

				state.StaticIPConfiguration = flattenSimStaticIPPropertiesDataSource(prop.StaticIPConfiguration)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func flattenSimStaticIPPropertiesDataSource(inputList *[]sim.SimStaticIPProperties) []SimStaticIPPropertiesDataSourceModel {
	outputList := make([]SimStaticIPPropertiesDataSourceModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := SimStaticIPPropertiesDataSourceModel{}

		if input.AttachedDataNetwork != nil {
			output.AttachedDataNetworkId = input.AttachedDataNetwork.Id
		}

		if input.Slice != nil {
			output.SliceId = input.Slice.Id
		}

		if input.StaticIP != nil && input.StaticIP.IPv4Address != nil {
			output.StaticIP = *input.StaticIP.IPv4Address
		}

		outputList = append(outputList, output)
	}

	return outputList
}
