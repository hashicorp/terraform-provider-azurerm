package network

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/routingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ManagerRoutingRuleResource{}

type ManagerRoutingRuleResource struct{}

func (ManagerRoutingRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return routingrules.ValidateRuleID
}

func (ManagerRoutingRuleResource) ResourceType() string {
	return "azurerm_network_manager_routing_rule"
}

func (ManagerRoutingRuleResource) ModelObject() interface{} {
	return &ManagerRoutingRuleResourceModel{}
}

type ManagerRoutingRuleResourceModel struct {
	Description      string                          `tfschema:"description"`
	Destination      []ManagerRoutingRuleDestination `tfschema:"destination"`
	Name             string                          `tfschema:"name"`
	NextHop          []ManagerRoutingRuleNextHop     `tfschema:"next_hop"`
	RuleCollectionId string                          `tfschema:"rule_collection_id"`
}

type ManagerRoutingRuleDestination struct {
	Address string `tfschema:"address"`
	Type    string `tfschema:"type"`
}

type ManagerRoutingRuleNextHop struct {
	Address string `tfschema:"address"`
	Type    string `tfschema:"type"`
}

func (ManagerRoutingRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"rule_collection_id": commonschema.ResourceIDReferenceRequiredForceNew(&routingrules.RuleCollectionId{}),

		"destination": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(routingrules.PossibleValuesForRoutingRuleDestinationType(), false),
					},
				},
			},
		},

		"next_hop": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(routingrules.PossibleValuesForRoutingRuleNextHopType(), false),
					},
					"address": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (ManagerRoutingRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerRoutingRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RoutingRules
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerRoutingRuleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			ruleCollectionId, err := routingrules.ParseRuleCollectionID(config.RuleCollectionId)
			if err != nil {
				return err
			}

			id := routingrules.NewRuleID(subscriptionId, ruleCollectionId.ResourceGroupName, ruleCollectionId.NetworkManagerName, ruleCollectionId.RoutingConfigurationName, ruleCollectionId.RuleCollectionName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			nextHop, err := expandNetworkManagerRoutingRuleNextHop(config.NextHop)
			if err != nil {
				return fmt.Errorf("expanding `next_hop`: %+v", err)
			}

			payload := routingrules.RoutingRule{
				Name: pointer.To(config.Name),
				Properties: &routingrules.RoutingRulePropertiesFormat{
					Description: pointer.To(config.Description),
					Destination: expandNetworkManagerRoutingRuleDestination(config.Destination),
					NextHop:     *nextHop,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ManagerRoutingRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RoutingRules

			id, err := routingrules.ParseRuleID(metadata.ResourceData.Id())
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

			schema := ManagerRoutingRuleResourceModel{
				Name:             id.RuleName,
				RuleCollectionId: routingrules.NewRuleCollectionID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.RoutingConfigurationName, id.RuleCollectionName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.Description = pointer.From(model.Properties.Description)
					schema.Destination = flattenNetworkManagerRoutingRuleDestination(props.Destination)
					schema.NextHop = flattenNetworkManagerRoutingRuleNextHop(props.NextHop)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ManagerRoutingRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RoutingRules

			id, err := routingrules.ParseRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerRoutingRuleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			parameters := resp.Model

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("destination") {
				parameters.Properties.Destination = expandNetworkManagerRoutingRuleDestination(model.Destination)
			}

			if metadata.ResourceData.HasChange("next_hop") {
				nextHop, err := expandNetworkManagerRoutingRuleNextHop(model.NextHop)
				if err != nil {
					return fmt.Errorf("expanding `next_hop`: %+v", err)
				}

				parameters.Properties.NextHop = *nextHop
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r ManagerRoutingRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RoutingRules

			id, err := routingrules.ParseRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, routingrules.DeleteOperationOptions{
				Force: pointer.To(true),
			}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandNetworkManagerRoutingRuleDestination(input []ManagerRoutingRuleDestination) routingrules.RoutingRuleRouteDestination {
	if len(input) == 0 {
		return routingrules.RoutingRuleRouteDestination{}
	}

	v := input[0]
	return routingrules.RoutingRuleRouteDestination{
		DestinationAddress: v.Address,
		Type:               routingrules.RoutingRuleDestinationType(v.Type),
	}
}

func flattenNetworkManagerRoutingRuleDestination(input routingrules.RoutingRuleRouteDestination) []ManagerRoutingRuleDestination {
	return []ManagerRoutingRuleDestination{
		{
			Address: input.DestinationAddress,
			Type:    string(input.Type),
		},
	}
}

func expandNetworkManagerRoutingRuleNextHop(input []ManagerRoutingRuleNextHop) (*routingrules.RoutingRuleNextHop, error) {
	if len(input) == 0 {
		return &routingrules.RoutingRuleNextHop{}, nil
	}

	v := input[0]

	if strings.EqualFold(v.Type, string(routingrules.RoutingRuleNextHopTypeVirtualAppliance)) && v.Address == "" {
		return nil, fmt.Errorf("address is required when type is `VirtualAppliance`")
	}

	return &routingrules.RoutingRuleNextHop{
		NextHopAddress: pointer.To(v.Address),
		NextHopType:    routingrules.RoutingRuleNextHopType(v.Type),
	}, nil
}

func flattenNetworkManagerRoutingRuleNextHop(input routingrules.RoutingRuleNextHop) []ManagerRoutingRuleNextHop {
	return []ManagerRoutingRuleNextHop{
		{
			Address: pointer.From(input.NextHopAddress),
			Type:    string(input.NextHopType),
		},
	}
}
