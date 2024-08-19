// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RouteMapModel struct {
	Name         string `tfschema:"name"`
	VirtualHubId string `tfschema:"virtual_hub_id"`
	Rules        []Rule `tfschema:"rule"`
}

type Rule struct {
	Actions           []Action             `tfschema:"action"`
	MatchCriteria     []Criterion          `tfschema:"match_criterion"`
	Name              string               `tfschema:"name"`
	NextStepIfMatched virtualwans.NextStep `tfschema:"next_step_if_matched"`
}

type Action struct {
	Parameters []Parameter                    `tfschema:"parameter"`
	Type       virtualwans.RouteMapActionType `tfschema:"type"`
}

type Parameter struct {
	AsPath      []string `tfschema:"as_path"`
	Community   []string `tfschema:"community"`
	RoutePrefix []string `tfschema:"route_prefix"`
}

type Criterion struct {
	AsPath         []string                           `tfschema:"as_path"`
	Community      []string                           `tfschema:"community"`
	MatchCondition virtualwans.RouteMapMatchCondition `tfschema:"match_condition"`
	RoutePrefix    []string                           `tfschema:"route_prefix"`
}

type RouteMapResource struct{}

var _ sdk.ResourceWithUpdate = RouteMapResource{}
var _ sdk.ResourceWithCustomizeDiff = RouteMapResource{}

func (r RouteMapResource) ResourceType() string {
	return "azurerm_route_map"
}

func (r RouteMapResource) ModelObject() interface{} {
	return &RouteMapModel{}
}

func (r RouteMapResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualwans.ValidateRouteMapID
}

func (r RouteMapResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RouteMapName,
		},

		"virtual_hub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualwans.ValidateVirtualHubID,
		},

		"rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"parameter": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"as_path": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},

											"community": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},

											"route_prefix": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
								},

								"type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(virtualwans.RouteMapActionTypeAdd),
										string(virtualwans.RouteMapActionTypeDrop),
										string(virtualwans.RouteMapActionTypeRemove),
										string(virtualwans.RouteMapActionTypeReplace),
										string(virtualwans.RouteMapActionTypeUnknown),
									}, false),
								},
							},
						},
					},

					"match_criterion": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"match_condition": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(virtualwans.RouteMapMatchConditionContains),
										string(virtualwans.RouteMapMatchConditionEquals),
										string(virtualwans.RouteMapMatchConditionNotContains),
										string(virtualwans.RouteMapMatchConditionNotEquals),
										string(virtualwans.RouteMapMatchConditionUnknown),
									}, false),
								},

								"as_path": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},

								"community": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},

								"route_prefix": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},

					"next_step_if_matched": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(virtualwans.NextStepUnknown),
						ValidateFunc: validation.StringInSlice([]string{
							string(virtualwans.NextStepContinue),
							string(virtualwans.NextStepTerminate),
							string(virtualwans.NextStepUnknown),
						}, false),
					},
				},
			},
		},
	}
}

func (r RouteMapResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r RouteMapResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RouteMapModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.VirtualWANs
			virtualHubId, err := virtualwans.ParseVirtualHubID(model.VirtualHubId)
			if err != nil {
				return err
			}

			id := virtualwans.NewRouteMapID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroupName, virtualHubId.VirtualHubName, model.Name)
			existing, err := client.RouteMapsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := &virtualwans.RouteMap{
				Properties: &virtualwans.RouteMapProperties{
					Rules: expandRules(model.Rules),
				},
			}

			if err := client.RouteMapsCreateOrUpdateThenPoll(ctx, id, *props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r RouteMapResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualWANs

			id, err := virtualwans.ParseRouteMapID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model RouteMapModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.RouteMapsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			if metadata.ResourceData.HasChange("rule") {
				existing.Model.Properties.Rules = expandRules(model.Rules)
			}

			if err := client.RouteMapsCreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r RouteMapResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualWANs

			id, err := virtualwans.ParseRouteMapID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.RouteMapsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := RouteMapModel{
				Name:         id.RouteMapName,
				VirtualHubId: virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Rules = flattenRules(props.Rules)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RouteMapResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualWANs

			id, err := virtualwans.ParseRouteMapID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.RouteMapsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r RouteMapResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config RouteMapModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Validate that all actions have parameters when they are not of type Drop
			for _, rule := range config.Rules {
				for _, action := range rule.Actions {
					if action.Type != virtualwans.RouteMapActionTypeDrop && len(action.Parameters) == 0 {
						return fmt.Errorf("parameters are required when rule action is not `Drop`")
					}
				}
			}

			return nil
		},
	}
}

func expandRules(input []Rule) *[]virtualwans.RouteMapRule {
	var rules []virtualwans.RouteMapRule
	if input == nil {
		return nil
	}

	for _, v := range input {
		rule := virtualwans.RouteMapRule{
			Name:          pointer.To(v.Name),
			Actions:       expandActions(v.Actions),
			MatchCriteria: expandCriteria(v.MatchCriteria),
		}

		if v.NextStepIfMatched != "" {
			rule.NextStepIfMatched = pointer.To(v.NextStepIfMatched)
		}

		rules = append(rules, rule)
	}

	return &rules
}

func expandActions(input []Action) *[]virtualwans.Action {
	var actions []virtualwans.Action
	if input == nil {
		return nil
	}

	for _, v := range input {
		action := virtualwans.Action{
			Type:       pointer.To(v.Type),
			Parameters: expandParameters(v.Parameters),
		}

		actions = append(actions, action)
	}

	return &actions
}

func expandParameters(input []Parameter) *[]virtualwans.Parameter {
	var parameters []virtualwans.Parameter
	if input == nil {
		return nil
	}

	for _, item := range input {
		v := item
		parameter := virtualwans.Parameter{}

		if v.AsPath != nil {
			parameter.AsPath = &v.AsPath
		}

		if v.Community != nil {
			parameter.Community = &v.Community
		}

		if v.RoutePrefix != nil {
			parameter.RoutePrefix = &v.RoutePrefix
		}

		parameters = append(parameters, parameter)
	}

	return &parameters
}

func expandCriteria(input []Criterion) *[]virtualwans.Criterion {
	var criteria []virtualwans.Criterion
	if input == nil {
		return nil
	}

	for _, item := range input {
		v := item
		criterion := virtualwans.Criterion{
			MatchCondition: pointer.To(v.MatchCondition),
		}

		if v.AsPath != nil {
			criterion.AsPath = &v.AsPath
		}

		if v.Community != nil {
			criterion.Community = &v.Community
		}

		if v.RoutePrefix != nil {
			criterion.RoutePrefix = &v.RoutePrefix
		}

		criteria = append(criteria, criterion)
	}

	return &criteria
}

func flattenRules(input *[]virtualwans.RouteMapRule) []Rule {
	var rules []Rule
	if input == nil {
		return rules
	}

	for _, v := range *input {
		rule := Rule{
			Actions:       flattenActions(v.Actions),
			MatchCriteria: flattenCriteria(v.MatchCriteria),
		}

		if v.Name != nil {
			rule.Name = *v.Name
		}

		if v.NextStepIfMatched != nil {
			rule.NextStepIfMatched = pointer.From(v.NextStepIfMatched)
		}

		rules = append(rules, rule)
	}

	return rules
}

func flattenActions(input *[]virtualwans.Action) []Action {
	var actions []Action
	if input == nil {
		return actions
	}

	for _, v := range *input {
		action := Action{
			Parameters: flattenParameters(v.Parameters),
		}

		if v.Type != nil {
			action.Type = pointer.From(v.Type)
		}

		actions = append(actions, action)
	}

	return actions
}

func flattenParameters(input *[]virtualwans.Parameter) []Parameter {
	var parameters []Parameter
	if input == nil {
		return parameters
	}

	for _, v := range *input {
		parameter := Parameter{}

		if v.AsPath != nil {
			parameter.AsPath = *v.AsPath
		}

		if v.Community != nil {
			parameter.Community = *v.Community
		}

		if v.RoutePrefix != nil {
			parameter.RoutePrefix = *v.RoutePrefix
		}

		parameters = append(parameters, parameter)
	}

	return parameters
}

func flattenCriteria(input *[]virtualwans.Criterion) []Criterion {
	var criteria []Criterion
	if input == nil {
		return criteria
	}

	for _, v := range *input {
		criterion := Criterion{}

		if v.AsPath != nil {
			criterion.AsPath = *v.AsPath
		}

		if v.Community != nil {
			criterion.Community = *v.Community
		}

		if v.MatchCondition != nil {
			criterion.MatchCondition = pointer.From(v.MatchCondition)
		}

		if v.RoutePrefix != nil {
			criterion.RoutePrefix = *v.RoutePrefix
		}

		criteria = append(criteria, criterion)
	}

	return criteria
}
