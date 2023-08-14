// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/adminrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerAdminRuleModel struct {
	Name                    string                                        `tfschema:"name"`
	NetworkRuleCollectionId string                                        `tfschema:"admin_rule_collection_id"`
	Action                  adminrules.SecurityConfigurationRuleAccess    `tfschema:"action"`
	Description             string                                        `tfschema:"description"`
	DestinationPortRanges   []string                                      `tfschema:"destination_port_ranges"`
	Destinations            []AddressPrefixItemModel                      `tfschema:"destination"`
	Direction               adminrules.SecurityConfigurationRuleDirection `tfschema:"direction"`
	Priority                int64                                         `tfschema:"priority"`
	Protocol                adminrules.SecurityConfigurationRuleProtocol  `tfschema:"protocol"`
	SourcePortRanges        []string                                      `tfschema:"source_port_ranges"`
	Sources                 []AddressPrefixItemModel                      `tfschema:"source"`
}

type AddressPrefixItemModel struct {
	AddressPrefix     string                       `tfschema:"address_prefix"`
	AddressPrefixType adminrules.AddressPrefixType `tfschema:"address_prefix_type"`
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
	return adminrules.ValidateRuleID
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
			ValidateFunc: adminrules.ValidateRuleCollectionID,
		},

		"action": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(adminrules.SecurityConfigurationRuleAccessAllow),
				string(adminrules.SecurityConfigurationRuleAccessDeny),
				string(adminrules.SecurityConfigurationRuleAccessAlwaysAllow),
			}, false),
		},

		"direction": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(adminrules.SecurityConfigurationRuleDirectionInbound),
				string(adminrules.SecurityConfigurationRuleDirectionOutbound),
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
				string(adminrules.SecurityConfigurationRuleProtocolAh),
				string(adminrules.SecurityConfigurationRuleProtocolAny),
				string(adminrules.SecurityConfigurationRuleProtocolIcmp),
				string(adminrules.SecurityConfigurationRuleProtocolEsp),
				string(adminrules.SecurityConfigurationRuleProtocolTcp),
				string(adminrules.SecurityConfigurationRuleProtocolUdp),
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
							string(adminrules.AddressPrefixTypeIPPrefix),
							string(adminrules.AddressPrefixTypeServiceTag),
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
							string(adminrules.AddressPrefixTypeIPPrefix),
							string(adminrules.AddressPrefixTypeServiceTag),
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

			client := metadata.Client.Network.AdminRules
			ruleCollectionId, err := adminrules.ParseRuleCollectionID(model.NetworkRuleCollectionId)
			if err != nil {
				return err
			}

			id := adminrules.NewRuleID(ruleCollectionId.SubscriptionId, ruleCollectionId.ResourceGroupName,
				ruleCollectionId.NetworkManagerName, ruleCollectionId.SecurityAdminConfigurationName, ruleCollectionId.RuleCollectionName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			rule := adminrules.AdminRule{
				Properties: &adminrules.AdminPropertiesFormat{
					Access:                model.Action,
					Destinations:          expandAddressPrefixItemModel(model.Destinations),
					DestinationPortRanges: &model.DestinationPortRanges,
					Direction:             model.Direction,
					Priority:              model.Priority,
					Protocol:              model.Protocol,
					SourcePortRanges:      &model.SourcePortRanges,
					Sources:               expandAddressPrefixItemModel(model.Sources),
				},
			}

			if model.Description != "" {
				rule.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, rule); err != nil {
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
			client := metadata.Client.Network.AdminRules

			id, err := adminrules.ParseRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerAdminRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			var rule adminrules.AdminRule
			if adminRule, ok := (*existing.Model).(adminrules.AdminRule); ok {
				rule = adminRule
			}

			if rule.Properties == nil {
				return fmt.Errorf("retrieving %s: property was nil", *id)
			}

			properties := rule.Properties

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
				properties.Priority = model.Priority
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

			if _, err := client.CreateOrUpdate(ctx, *id, rule); err != nil {
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
			client := metadata.Client.Network.AdminRules

			id, err := adminrules.ParseRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			var rule adminrules.AdminRule
			if adminRule, ok := (*existing.Model).(adminrules.AdminRule); ok {
				rule = adminRule
			}

			if rule.Properties == nil {
				return fmt.Errorf("retrieving %s: property was nil", *id)
			}

			properties := rule.Properties

			state := ManagerAdminRuleModel{
				Action: properties.Access,
				Name:   id.RuleName,
				NetworkRuleCollectionId: adminrules.NewRuleCollectionID(id.SubscriptionId, id.ResourceGroupName,
					id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName).ID(),
				Destinations: flattenAddressPrefixItemModel(properties.Destinations),
				Direction:    properties.Direction,
				Priority:     properties.Priority,
				Protocol:     properties.Protocol,
				Sources:      flattenAddressPrefixItemModel(properties.Sources),
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.DestinationPortRanges != nil {
				state.DestinationPortRanges = *properties.DestinationPortRanges
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
			client := metadata.Client.Network.AdminRules

			id, err := adminrules.ParseRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id, adminrules.DeleteOperationOptions{
				Force: utils.Bool(true),
			})
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandAddressPrefixItemModel(inputList []AddressPrefixItemModel) *[]adminrules.AddressPrefixItem {
	var outputList []adminrules.AddressPrefixItem
	for _, v := range inputList {
		input := v
		output := adminrules.AddressPrefixItem{
			AddressPrefixType: &input.AddressPrefixType,
		}

		if input.AddressPrefix != "" {
			output.AddressPrefix = &input.AddressPrefix
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenAddressPrefixItemModel(inputList *[]adminrules.AddressPrefixItem) []AddressPrefixItemModel {
	var outputList []AddressPrefixItemModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := AddressPrefixItemModel{}

		if input.AddressPrefixType != nil {
			output.AddressPrefixType = *input.AddressPrefixType
		}

		if input.AddressPrefix != nil {
			output.AddressPrefix = *input.AddressPrefix
		}

		outputList = append(outputList, output)
	}

	return outputList
}
