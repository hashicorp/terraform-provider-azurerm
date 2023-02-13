package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type ManagerAdminRuleModel struct {
	Name                    string                                     `tfschema:"name"`
	NetworkRuleCollectionId string                                     `tfschema:"admin_rule_collection_id"`
	Action                  network.SecurityConfigurationRuleAccess    `tfschema:"action"`
	Description             string                                     `tfschema:"description"`
	DestinationPortRanges   []string                                   `tfschema:"destination_port_ranges"`
	Destinations            []AddressPrefixItemModel                   `tfschema:"destination"`
	Direction               network.SecurityConfigurationRuleDirection `tfschema:"direction"`
	Priority                int32                                      `tfschema:"priority"`
	Protocol                network.SecurityConfigurationRuleProtocol  `tfschema:"protocol"`
	SourcePortRanges        []string                                   `tfschema:"source_port_ranges"`
	Sources                 []AddressPrefixItemModel                   `tfschema:"source"`
}

type AddressPrefixItemModel struct {
	AddressPrefix     string                    `tfschema:"address_prefix"`
	AddressPrefixType network.AddressPrefixType `tfschema:"address_prefix_type"`
}

type ManagerAdminRuleResource struct{}

var _ sdk.ResourceWithUpdate = ManagerAdminRuleResource{}

func (r ManagerAdminRuleResource) ResourceType() string {
	return "azurerm_network_manager_admin_rule"
}

func (r ManagerAdminRuleResource) ModelObject() interface{} {
	return &ManagerAdminRuleModel{}
}

func (r ManagerAdminRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerAdminRuleID
}

func (r ManagerAdminRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"admin_rule_collection_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NetworkManagerAdminRuleCollectionID,
		},

		"action": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(network.SecurityConfigurationRuleAccessAllow),
				string(network.SecurityConfigurationRuleAccessDeny),
				string(network.SecurityConfigurationRuleAccessAlwaysAllow),
			}, false),
		},

		"direction": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(network.SecurityConfigurationRuleDirectionInbound),
				string(network.SecurityConfigurationRuleDirectionOutbound),
			}, false),
		},

		"priority": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 4096),
		},

		"protocol": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(network.SecurityConfigurationRuleProtocolAh),
				string(network.SecurityConfigurationRuleProtocolAny),
				string(network.SecurityConfigurationRuleProtocolIcmp),
				string(network.SecurityConfigurationRuleProtocolEsp),
				string(network.SecurityConfigurationRuleProtocolTCP),
				string(network.SecurityConfigurationRuleProtocolUDP),
			}, false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"destination_port_ranges": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"destination": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_prefix": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_prefix_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(network.AddressPrefixTypeIPPrefix),
							string(network.AddressPrefixTypeServiceTag),
						}, false),
					},
				},
			},
		},

		"source_port_ranges": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"source": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_prefix": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_prefix_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(network.AddressPrefixTypeIPPrefix),
							string(network.AddressPrefixTypeServiceTag),
						}, false),
					},
				},
			},
		},
	}
}

func (r ManagerAdminRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerAdminRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerAdminRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ManagerAdminRulesClient
			ruleCollectionId, err := parse.NetworkManagerAdminRuleCollectionID(model.NetworkRuleCollectionId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerAdminRuleID(ruleCollectionId.SubscriptionId, ruleCollectionId.ResourceGroup,
				ruleCollectionId.NetworkManagerName, ruleCollectionId.SecurityAdminConfigurationName, ruleCollectionId.RuleCollectionName, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName, id.RuleName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			rule := &network.AdminRule{
				AdminPropertiesFormat: &network.AdminPropertiesFormat{
					Access:                model.Action,
					Destinations:          expandAddressPrefixItemModel(model.Destinations),
					DestinationPortRanges: &model.DestinationPortRanges,
					Direction:             model.Direction,
					Priority:              utils.Int32(model.Priority),
					Protocol:              model.Protocol,
					SourcePortRanges:      &model.SourcePortRanges,
					Sources:               expandAddressPrefixItemModel(model.Sources),
				},
			}

			if model.Description != "" {
				rule.AdminPropertiesFormat.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, *rule, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName, id.RuleName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerAdminRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerAdminRulesClient

			id, err := parse.NetworkManagerAdminRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerAdminRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName, id.RuleName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var rule *network.AdminRule
			if adminRule, ok := existing.Value.AsAdminRule(); ok {
				rule = adminRule
			}

			properties := rule.AdminPropertiesFormat
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("action") {
				properties.Access = model.Action
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Description = &model.Description
				} else {
					properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("destination_port_ranges") {
				properties.DestinationPortRanges = &model.DestinationPortRanges
			}

			if metadata.ResourceData.HasChange("destination") {
				properties.Destinations = expandAddressPrefixItemModel(model.Destinations)
			}

			if metadata.ResourceData.HasChange("direction") {
				properties.Direction = model.Direction
			}

			if metadata.ResourceData.HasChange("priority") {
				properties.Priority = utils.Int32(model.Priority)
			}

			if metadata.ResourceData.HasChange("protocol") {
				properties.Protocol = model.Protocol
			}

			if metadata.ResourceData.HasChange("source_port_ranges") {
				properties.SourcePortRanges = &model.SourcePortRanges
			}

			if metadata.ResourceData.HasChange("source") {
				properties.Sources = expandAddressPrefixItemModel(model.Sources)
			}

			if _, err := client.CreateOrUpdate(ctx, rule, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName, id.RuleName); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerAdminRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerAdminRulesClient

			id, err := parse.NetworkManagerAdminRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName, id.RuleName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var rule *network.AdminRule
			if adminRule, ok := existing.Value.AsAdminRule(); ok {
				rule = adminRule
			}

			properties := rule.AdminPropertiesFormat
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ManagerAdminRuleModel{
				Action: properties.Access,
				Name:   id.RuleName,
				NetworkRuleCollectionId: parse.NewNetworkManagerAdminRuleCollectionID(id.SubscriptionId, id.ResourceGroup,
					id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName).ID(),
				Destinations: flattenAddressPrefixItemModel(properties.Destinations),
				Direction:    properties.Direction,
				Protocol:     properties.Protocol,
				Sources:      flattenAddressPrefixItemModel(properties.Sources),
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.DestinationPortRanges != nil {
				state.DestinationPortRanges = *properties.DestinationPortRanges
			}

			state.Priority = 0
			if properties.Priority != nil {
				state.Priority = *properties.Priority
			}

			if properties.SourcePortRanges != nil {
				state.SourcePortRanges = *properties.SourcePortRanges
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerAdminRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerAdminRulesClient

			id, err := parse.NetworkManagerAdminRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName, id.RuleName, utils.Bool(true))
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandAddressPrefixItemModel(inputList []AddressPrefixItemModel) *[]network.AddressPrefixItem {
	var outputList []network.AddressPrefixItem
	for _, v := range inputList {
		input := v
		output := network.AddressPrefixItem{
			AddressPrefixType: input.AddressPrefixType,
		}

		if input.AddressPrefix != "" {
			output.AddressPrefix = &input.AddressPrefix
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenAddressPrefixItemModel(inputList *[]network.AddressPrefixItem) []AddressPrefixItemModel {
	var outputList []AddressPrefixItemModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := AddressPrefixItemModel{
			AddressPrefixType: input.AddressPrefixType,
		}

		if input.AddressPrefix != nil {
			output.AddressPrefix = *input.AddressPrefix
		}

		outputList = append(outputList, output)
	}

	return outputList
}
