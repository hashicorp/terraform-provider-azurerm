// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddoscustompolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity

type NetworkDDoSCustomPolicyResource struct{}

var (
	_ sdk.Resource                  = NetworkDDoSCustomPolicyResource{}
	_ sdk.ResourceWithIdentity      = NetworkDDoSCustomPolicyResource{}
	_ sdk.ResourceWithUpdate        = NetworkDDoSCustomPolicyResource{}
	_ sdk.ResourceWithCustomizeDiff = NetworkDDoSCustomPolicyResource{}
)

type NetworkDDoSCustomPolicyResourceModel struct {
	Name              string                                      `tfschema:"name"`
	ResourceGroupName string                                      `tfschema:"resource_group_name"`
	Location          string                                      `tfschema:"location"`
	DetectionRules    []NetworkDDoSCustomPolicyDetectionRuleModel `tfschema:"detection_rule"`
	Tags              map[string]string                           `tfschema:"tags"`

	PublicIPAddressIds []string `tfschema:"public_ip_address_ids"`
}

type NetworkDDoSCustomPolicyDetectionRuleModel struct {
	Name             string `tfschema:"name"`
	PacketsPerSecond int64  `tfschema:"packets_per_second"`
	TrafficType      string `tfschema:"traffic_type"`
}

func (NetworkDDoSCustomPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"detection_rule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 3, // one rule per traffic type: Tcp, Udp, TcpSyn
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"packets_per_second": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						// NOTE: the allowed range differs per `traffic_type`, so this is a coarse
						// guard covering the union of them - the exact range is checked in CustomizeDiff.
						ValidateFunc: validation.IntBetween(1000, 2000000),
					},
					"traffic_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(ddoscustompolicies.PossibleValuesForDdosTrafficType(), false),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (NetworkDDoSCustomPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"public_ip_address_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (NetworkDDoSCustomPolicyResource) ModelObject() interface{} {
	return &NetworkDDoSCustomPolicyResourceModel{}
}

func (r NetworkDDoSCustomPolicyResource) ResourceType() string {
	return "azurerm_network_ddos_custom_policy"
}

func (r NetworkDDoSCustomPolicyResource) Identity() resourceids.ResourceId {
	return &ddoscustompolicies.DdosCustomPolicyId{}
}

func (NetworkDDoSCustomPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return ddoscustompolicies.ValidateDdosCustomPolicyID
}

func (r NetworkDDoSCustomPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.DdosCustomPoliciesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config NetworkDDoSCustomPolicyResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := ddoscustompolicies.NewDdosCustomPolicyID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := ddoscustompolicies.DdosCustomPolicy{
				Location: pointer.To(location.Normalize(config.Location)),
				Properties: &ddoscustompolicies.DdosCustomPolicyPropertiesFormat{
					DetectionRules: expandNetworkDDoSCustomPolicyDetectionRules(config.DetectionRules),
				},
				Tags: pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandNetworkDDoSCustomPolicyDetectionRules(input []NetworkDDoSCustomPolicyDetectionRuleModel) *[]ddoscustompolicies.DdosDetectionRule {
	result := make([]ddoscustompolicies.DdosDetectionRule, 0, len(input))
	for _, v := range input {
		result = append(result, ddoscustompolicies.DdosDetectionRule{
			Name: pointer.To(v.Name),
			Properties: &ddoscustompolicies.DdosDetectionRulePropertiesFormat{
				// NOTE: `detectionMode` isn't exposed in the schema since `DdosDetectionMode` only has one single possible value.
				DetectionMode: pointer.To(ddoscustompolicies.DdosDetectionModeTrafficThreshold),
				TrafficDetectionRule: &ddoscustompolicies.TrafficDetectionRule{
					PacketsPerSecond: pointer.To(v.PacketsPerSecond),
					TrafficType:      pointer.ToEnum[ddoscustompolicies.DdosTrafficType](v.TrafficType),
				},
			},
		})
	}
	return &result
}

// packetsPerSecondRangeForTrafficType returns the allowed `packets_per_second` range for each
// `traffic_type`. These ranges are enforced by the service but aren't published in the API definition.
var packetsPerSecondRangeForTrafficType = map[string]struct {
	min int64
	max int64
}{
	string(ddoscustompolicies.DdosTrafficTypeTcp):    {min: 50000, max: 2000000},
	string(ddoscustompolicies.DdosTrafficTypeUdp):    {min: 20000, max: 200000},
	string(ddoscustompolicies.DdosTrafficTypeTcpSyn): {min: 1000, max: 100000},
}

func (r NetworkDDoSCustomPolicyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config NetworkDDoSCustomPolicyResourceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			names := make(map[string]struct{}, len(config.DetectionRules))
			trafficTypes := make(map[string]struct{}, len(config.DetectionRules))

			for _, rule := range config.DetectionRules {
				if _, exists := names[rule.Name]; exists {
					return fmt.Errorf("the `detection_rule` blocks must have unique `name` values, got duplicate `%s`", rule.Name)
				}
				names[rule.Name] = struct{}{}

				if _, exists := trafficTypes[rule.TrafficType]; exists {
					return fmt.Errorf("the `detection_rule` blocks must have unique `traffic_type` values, got duplicate `%s`", rule.TrafficType)
				}
				trafficTypes[rule.TrafficType] = struct{}{}

				if allowed, ok := packetsPerSecondRangeForTrafficType[rule.TrafficType]; ok {
					if rule.PacketsPerSecond < allowed.min || rule.PacketsPerSecond > allowed.max {
						return fmt.Errorf("`packets_per_second` must be between %d and %d when `traffic_type` is `%s`, got %d", allowed.min, allowed.max, rule.TrafficType, rule.PacketsPerSecond)
					}
				}
			}

			return nil
		},
	}
}

func (r NetworkDDoSCustomPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.DdosCustomPoliciesClient

			id, err := ddoscustompolicies.ParseDdosCustomPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			return r.flatten(metadata, id, model)
		},
	}
}

func (r NetworkDDoSCustomPolicyResource) flatten(metadata sdk.ResourceMetaData, id *ddoscustompolicies.DdosCustomPolicyId, model *ddoscustompolicies.DdosCustomPolicy) error {
	state := NetworkDDoSCustomPolicyResourceModel{
		Name:              id.DdosCustomPolicyName,
		ResourceGroupName: id.ResourceGroupName,
		Location:          location.NormalizeNilable(model.Location),
		Tags:              pointer.From(model.Tags),
	}

	if props := model.Properties; props != nil {
		state.DetectionRules = flattenNetworkDDoSCustomPolicyDetectionRules(props.DetectionRules)

		publicIPAddressIds, err := flattenNetworkDDoSCustomPolicyPublicIPAddressIds(props.PublicIPAddresses)
		if err != nil {
			return err
		}
		state.PublicIPAddressIds = publicIPAddressIds
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenNetworkDDoSCustomPolicyDetectionRules(input *[]ddoscustompolicies.DdosDetectionRule) []NetworkDDoSCustomPolicyDetectionRuleModel {
	output := make([]NetworkDDoSCustomPolicyDetectionRuleModel, 0)
	if input == nil {
		return output
	}

	for _, v := range *input {
		rule := NetworkDDoSCustomPolicyDetectionRuleModel{
			Name: pointer.From(v.Name),
		}

		// NOTE: `detectionMode` isn't exposed in the schema, so it's intentionally not read back here.
		if props := v.Properties; props != nil {
			if trafficDetectionRule := props.TrafficDetectionRule; trafficDetectionRule != nil {
				rule.PacketsPerSecond = pointer.From(trafficDetectionRule.PacketsPerSecond)
				rule.TrafficType = string(pointer.From(trafficDetectionRule.TrafficType))
			}
		}

		output = append(output, rule)
	}

	return output
}

func flattenNetworkDDoSCustomPolicyPublicIPAddressIds(input *[]ddoscustompolicies.SubResource) ([]string, error) {
	output := make([]string, 0)
	if input == nil {
		return output, nil
	}

	for _, v := range *input {
		// NOTE: the casing of the IDs returned by the API isn't guaranteed to match the casing used
		// when the association was created, so these are normalised to prevent a diff.
		id, err := commonids.ParsePublicIPAddressIDInsensitively(pointer.From(v.Id))
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", pointer.From(v.Id), err)
		}

		output = append(output, id.ID())
	}

	return output, nil
}

func (r NetworkDDoSCustomPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.DdosCustomPoliciesClient

			id, err := ddoscustompolicies.ParseDdosCustomPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config NetworkDDoSCustomPolicyResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("detection_rule") {
				payload.Properties.DetectionRules = expandNetworkDDoSCustomPolicyDetectionRules(config.DetectionRules)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r NetworkDDoSCustomPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.DdosCustomPoliciesClient

			id, err := ddoscustompolicies.ParseDdosCustomPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
