package network

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/routingrulecollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = ManagerRoutingRuleCollectionResource{}

type ManagerRoutingRuleCollectionResource struct{}

func (ManagerRoutingRuleCollectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return routingrulecollections.ValidateRuleCollectionID
}

func (ManagerRoutingRuleCollectionResource) ResourceType() string {
	return "azurerm_network_manager_routing_rule_collection"
}

func (ManagerRoutingRuleCollectionResource) ModelObject() interface{} {
	return &ManagerRoutingRuleCollectionResourceModel{}
}

type ManagerRoutingRuleCollectionResourceModel struct {
	BgpRoutePropagationEnabled bool     `tfschema:"bgp_route_propagation_enabled"`
	Description                string   `tfschema:"description"`
	Name                       string   `tfschema:"name"`
	NetworkGroupIds            []string `tfschema:"network_group_ids"`
	RoutingConfigurationId     string   `tfschema:"routing_configuration_id"`
}

func (ManagerRoutingRuleCollectionResource) Arguments() map[string]*pluginsdk.Schema {
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

		"routing_configuration_id": commonschema.ResourceIDReferenceRequiredForceNew(&routingrulecollections.RoutingConfigurationId{}),

		"network_group_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem:     commonschema.ResourceIDReferenceElem(&networkgroups.NetworkGroupId{}),
		},

		"bgp_route_propagation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (ManagerRoutingRuleCollectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerRoutingRuleCollectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RoutingRuleCollections
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerRoutingRuleCollectionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			routingConfigurationId, err := routingrulecollections.ParseRoutingConfigurationID(config.RoutingConfigurationId)
			if err != nil {
				return err
			}

			id := routingrulecollections.NewRuleCollectionID(subscriptionId, routingConfigurationId.ResourceGroupName, routingConfigurationId.NetworkManagerName, routingConfigurationId.RoutingConfigurationName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := routingrulecollections.RoutingRuleCollection{
				Name: pointer.To(config.Name),
				Properties: &routingrulecollections.RoutingRuleCollectionPropertiesFormat{
					AppliesTo:                  expandNetworkManagerRoutingNetworkGroupIds(config.NetworkGroupIds),
					Description:                pointer.To(config.Description),
					DisableBgpRoutePropagation: expandBgpRoutePropagation(config.BgpRoutePropagationEnabled),
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

func (r ManagerRoutingRuleCollectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RoutingRuleCollections

			id, err := routingrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
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

			schema := ManagerRoutingRuleCollectionResourceModel{
				Name:                   id.RoutingConfigurationName,
				RoutingConfigurationId: routingrulecollections.NewRoutingConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.RoutingConfigurationName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.BgpRoutePropagationEnabled = flattenBgpRoutePropagation(props.DisableBgpRoutePropagation)
					schema.Description = pointer.From(props.Description)

					networkGroupIds, err := flattenNetworkManagerRoutingNetworkGroupIds(props.AppliesTo)
					if err != nil {
						return fmt.Errorf("flattening `network_group_ids`: %+v", err)
					}

					schema.NetworkGroupIds = networkGroupIds
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ManagerRoutingRuleCollectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RoutingRuleCollections

			id, err := routingrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerRoutingRuleCollectionResourceModel
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

			if metadata.ResourceData.HasChange("bgp_route_propagation_enabled") {
				parameters.Properties.DisableBgpRoutePropagation = expandBgpRoutePropagation(model.BgpRoutePropagationEnabled)
			}

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("network_group_ids") {
				parameters.Properties.AppliesTo = expandNetworkManagerRoutingNetworkGroupIds(model.NetworkGroupIds)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r ManagerRoutingRuleCollectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RoutingRuleCollections

			id, err := routingrulecollections.ParseRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, routingrulecollections.DeleteOperationOptions{
				Force: pointer.To(true),
			}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandBgpRoutePropagation(input bool) *routingrulecollections.DisableBgpRoutePropagation {
	if input {
		return pointer.To(routingrulecollections.DisableBgpRoutePropagationFalse)
	}
	return pointer.To(routingrulecollections.DisableBgpRoutePropagationTrue)
}

func flattenBgpRoutePropagation(input *routingrulecollections.DisableBgpRoutePropagation) bool {
	return pointer.From(input) == routingrulecollections.DisableBgpRoutePropagationFalse
}

func expandNetworkManagerRoutingNetworkGroupIds(inputList []string) []routingrulecollections.NetworkManagerRoutingGroupItem {
	outputList := make([]routingrulecollections.NetworkManagerRoutingGroupItem, 0, len(inputList))
	for _, v := range inputList {
		output := routingrulecollections.NetworkManagerRoutingGroupItem{
			NetworkGroupId: v,
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenNetworkManagerRoutingNetworkGroupIds(inputList []routingrulecollections.NetworkManagerRoutingGroupItem) ([]string, error) {
	outputList := make([]string, 0, len(inputList))

	for _, input := range inputList {
		networkGroupId, err := networkgroups.ParseNetworkGroupIDInsensitively(input.NetworkGroupId)
		if err != nil {
			return nil, err
		}

		outputList = append(outputList, networkGroupId.ID())
	}

	return outputList, nil
}
