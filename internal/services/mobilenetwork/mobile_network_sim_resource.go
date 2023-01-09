package mobilenetwork

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/attacheddatanetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/sim"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/slice"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SimModel struct {
	Name                                  string                       `tfschema:"name"`
	MobileNetworkSimGroupId               string                       `tfschema:"mobile_network_sim_group_id"`
	AuthenticationKey                     string                       `tfschema:"authentication_key"`
	DeviceType                            string                       `tfschema:"device_type"`
	IntegratedCircuitCardIdentifier       string                       `tfschema:"integrated_circuit_card_identifier"`
	InternationalMobileSubscriberIdentity string                       `tfschema:"international_mobile_subscriber_identity"`
	OperatorKeyCode                       string                       `tfschema:"operator_key_code"`
	SimPolicyId                           string                       `tfschema:"sim_policy_id"`
	StaticIPConfiguration                 []SimStaticIPPropertiesModel `tfschema:"static_ip_configuration"`
	SimState                              string                       `tfschema:"sim_state"`
	VendorKeyFingerprint                  string                       `tfschema:"vendor_key_fingerprint"`
	VendorName                            string                       `tfschema:"vendor_name"`
}

type SimStaticIPPropertiesModel struct {
	AttachedDataNetworkId string `tfschema:"attached_data_network_id"`
	SliceId               string `tfschema:"slice_id"`
	StaticIP              string `tfschema:"static_ipv4_address"`
}

type SimResource struct{}

var _ sdk.ResourceWithUpdate = SimResource{}

func (r SimResource) ResourceType() string {
	return "azurerm_mobile_network_sim"
}

func (r SimResource) ModelObject() interface{} {
	return &SimModel{}
}

func (r SimResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sim.ValidateSimID
}

func (r SimResource) Arguments() map[string]*pluginsdk.Schema {
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

		"authentication_key": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
			ForceNew:  true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9a-fA-F]{32}$`),
				"The authentication key must be a 32 character hexadecimal string.",
			),
		},

		"integrated_circuit_card_identifier": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^89[0-9]{17,18}$`),
				`The integrated circuit card ID (ICCID) must be a 19/20 digit number starts with "89".`,
			),
		},

		"international_mobile_subscriber_identity": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9]{15}$`),
				"The international mobile subscriber identity (IMSI) must be a 15 digit number.",
			),
		},

		"operator_key_code": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			ForceNew:  true,
			Sensitive: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9a-fA-F]{32}$`),
				"The operator key code (OPC) must be a 32 hexadecimal number.",
			),
		},

		"device_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sim_policy_id": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     simpolicy.ValidateSimPolicyID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"static_ip_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"attached_data_network_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: attacheddatanetwork.ValidateAttachedDataNetworkID,
					},

					"slice_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: slice.ValidateSliceID,
					},

					"static_ipv4_address": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.IPv4Address,
					},
				},
			},
		},
	}
}

func (r SimResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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

func (r SimResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SimModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SIMClient
			simGroupId, err := simgroup.ParseSimGroupID(model.MobileNetworkSimGroupId)
			if err != nil {
				return err
			}

			id := sim.NewSimID(simGroupId.SubscriptionId, simGroupId.ResourceGroupName, simGroupId.SimGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &sim.Sim{
				Properties: sim.SimPropertiesFormat{
					InternationalMobileSubscriberIdentity: model.InternationalMobileSubscriberIdentity,
				},
			}

			if model.SimPolicyId != "" {
				properties.Properties.SimPolicy = &sim.SimPolicyResourceId{
					Id: model.SimPolicyId,
				}
			}

			if model.AuthenticationKey != "" {
				properties.Properties.AuthenticationKey = &model.AuthenticationKey
			}

			if model.DeviceType != "" {
				properties.Properties.DeviceType = &model.DeviceType
			}

			if model.IntegratedCircuitCardIdentifier != "" {
				properties.Properties.IntegratedCircuitCardIdentifier = &model.IntegratedCircuitCardIdentifier
			}

			if model.OperatorKeyCode != "" {
				properties.Properties.OperatorKeyCode = &model.OperatorKeyCode
			}

			properties.Properties.StaticIPConfiguration = expandSimStaticIPPropertiesModel(model.StaticIPConfiguration)

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// AuthenticationKey and OperatorKeyCode are not returned from the API side.
			metadata.ResourceData.Set("authentication_key", model.AuthenticationKey)
			metadata.ResourceData.Set("operator_key_code", model.OperatorKeyCode)

			metadata.SetID(id)

			return nil
		},
	}
}

func (r SimResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SIMClient

			id, err := sim.ParseSimID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SimModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if model.AuthenticationKey != "" {
				properties.Properties.AuthenticationKey = &model.AuthenticationKey
			} else {
				properties.Properties.AuthenticationKey = nil
			}

			if metadata.ResourceData.HasChange("device_type") {
				if model.DeviceType != "" {
					properties.Properties.DeviceType = &model.DeviceType
				} else {
					properties.Properties.DeviceType = nil
				}
			}

			if metadata.ResourceData.HasChange("integrated_circuit_card_identifier") {
				if model.IntegratedCircuitCardIdentifier != "" {
					properties.Properties.IntegratedCircuitCardIdentifier = &model.IntegratedCircuitCardIdentifier
				} else {
					properties.Properties.IntegratedCircuitCardIdentifier = nil
				}
			}

			if model.OperatorKeyCode != "" {
				properties.Properties.OperatorKeyCode = &model.OperatorKeyCode
			} else {
				properties.Properties.OperatorKeyCode = nil
			}

			if metadata.ResourceData.HasChange("sim_policy") {
				properties.Properties.SimPolicy = &sim.SimPolicyResourceId{
					Id: model.SimPolicyId,
				}
			}

			if metadata.ResourceData.HasChange("static_ip_configuration") {
				properties.Properties.StaticIPConfiguration = expandSimStaticIPPropertiesModel(model.StaticIPConfiguration)
			}

			properties.SystemData = nil

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SimResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SIMClient

			id, err := sim.ParseSimID(metadata.ResourceData.Id())
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
				state.SimPolicyId = simPolicy.Id
			}

			if properties.SimState != nil {
				state.SimState = string(*properties.SimState)
			}

			state.StaticIPConfiguration = flattenSimStaticIPPropertiesModel(properties.StaticIPConfiguration)

			if properties.VendorKeyFingerprint != nil {
				state.VendorKeyFingerprint = *properties.VendorKeyFingerprint
			}

			if properties.VendorName != nil {
				state.VendorName = *properties.VendorName
			}

			state.AuthenticationKey = metadata.ResourceData.Get("authentication_key").(string)
			state.OperatorKeyCode = metadata.ResourceData.Get("operator_key_code").(string)

			return metadata.Encode(&state)
		},
	}
}

func (r SimResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SIMClient

			id, err := sim.ParseSimID(metadata.ResourceData.Id())
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

func expandSimStaticIPPropertiesModel(inputList []SimStaticIPPropertiesModel) *[]sim.SimStaticIPProperties {
	var outputList []sim.SimStaticIPProperties
	for _, v := range inputList {
		input := v
		output := sim.SimStaticIPProperties{
			AttachedDataNetwork: &sim.AttachedDataNetworkResourceId{
				Id: input.AttachedDataNetworkId,
			},
			Slice: &sim.SliceResourceId{
				Id: input.SliceId,
			},
			StaticIP: &sim.SimStaticIPPropertiesStaticIP{
				IPv4Address: &input.StaticIP,
			},
		}
		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenSimStaticIPPropertiesModel(inputList *[]sim.SimStaticIPProperties) []SimStaticIPPropertiesModel {
	var outputList []SimStaticIPPropertiesModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := SimStaticIPPropertiesModel{}

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
