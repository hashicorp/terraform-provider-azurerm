package mobilenetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/sim"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SimDataSource struct{}

var _ sdk.DataSource = SimDataSource{}

func (r SimDataSource) ResourceType() string {
	return "azurerm_mobile_network_sim"
}

func (r SimDataSource) ModelObject() interface{} {
	return &SimModel{}
}

func (r SimDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sim.ValidateSimID
}

func (r SimDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_sim_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
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
			var metaModel SimModel
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
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := SimModel{
				Name:                    id.SimName,
				MobileNetworkSimGroupId: simgroup.NewSimGroupID(id.SubscriptionId, id.ResourceGroupName, id.SimGroupName).ID(),
			}

			properties := &model.Properties

			if properties.DeviceType != nil {
				state.DeviceType = *properties.DeviceType
			}

			if properties.IntegratedCircuitCardIdentifier != nil {
				state.IntegratedCircuitCardIdentifier = *properties.IntegratedCircuitCardIdentifier
			}

			state.InternationalMobileSubscriberIdentity = properties.InternationalMobileSubscriberIdentity

			if simPolicy := properties.SimPolicy; properties.SimPolicy != nil {
				state.SimPolicyId = replaceUpperCaseWordsWorkAround(simPolicy.Id)
			}

			if properties.SimState != nil {
				state.SimState = string(*properties.SimState)
			}

			staticIPConfigurationValue, err := flattenSimStaticIPPropertiesModel(properties.StaticIPConfiguration)
			if err != nil {
				return err
			}
			state.StaticIPConfiguration = staticIPConfigurationValue

			if properties.VendorKeyFingerprint != nil {
				state.VendorKeyFingerprint = *properties.VendorKeyFingerprint
			}

			if properties.VendorName != nil {
				state.VendorName = *properties.VendorName
			}
			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
