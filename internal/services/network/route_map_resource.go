// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

type RouteMapModel struct {
	Name         string `tfschema:"name"`
	VirtualHubId string `tfschema:"virtual_hub_id"`
	Rules        []Rule `tfschema:"rule"`
}

type Rule struct {
	Actions           []Action         `tfschema:"action"`
	MatchCriteria     []Criterion      `tfschema:"match_criterion"`
	Name              string           `tfschema:"name"`
	NextStepIfMatched network.NextStep `tfschema:"next_step_if_matched"`
}

type Action struct {
	Parameters []Parameter                `tfschema:"parameter"`
	Type       network.RouteMapActionType `tfschema:"type"`
}

type Parameter struct {
	AsPath      []string `tfschema:"as_path"`
	Community   []string `tfschema:"community"`
	RoutePrefix []string `tfschema:"route_prefix"`
}

type Criterion struct {
	AsPath         []string                       `tfschema:"as_path"`
	Community      []string                       `tfschema:"community"`
	MatchCondition network.RouteMapMatchCondition `tfschema:"match_condition"`
	RoutePrefix    []string                       `tfschema:"route_prefix"`
}

type RouteMapResource struct{}

var _ sdk.ResourceWithUpdate = RouteMapResource{}

func (r RouteMapResource) ResourceType() string {
	return "azurerm_route_map"
}

func (r RouteMapResource) ModelObject() interface{} {
	return &RouteMapModel{}
}

func (r RouteMapResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.RouteMapID
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
			ValidateFunc: validate.VirtualHubID,
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
									Required: true,
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
										string(network.RouteMapActionTypeAdd),
										string(network.RouteMapActionTypeDrop),
										string(network.RouteMapActionTypeRemove),
										string(network.RouteMapActionTypeReplace),
										string(network.RouteMapActionTypeUnknown),
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
										string(network.RouteMapMatchConditionContains),
										string(network.RouteMapMatchConditionEquals),
										string(network.RouteMapMatchConditionNotContains),
										string(network.RouteMapMatchConditionNotEquals),
										string(network.RouteMapMatchConditionUnknown),
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
						Default:  string(network.NextStepUnknown),
						ValidateFunc: validation.StringInSlice([]string{
							string(network.NextStepContinue),
							string(network.NextStepTerminate),
							string(network.NextStepUnknown),
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
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RouteMapModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.RouteMapsClient
			virtualHubId, err := parse.VirtualHubID(model.VirtualHubId)
			if err != nil {
				return err
			}

			id := parse.NewRouteMapID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroup, virtualHubId.Name, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := &network.RouteMap{
				RouteMapProperties: &network.RouteMapProperties{
					Rules: expandRules(model.Rules),
				},
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualHubName, id.Name, *props)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r RouteMapResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RouteMapsClient

			id, err := parse.RouteMapID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model RouteMapModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("rule") {
				existing.RouteMapProperties.Rules = expandRules(model.Rules)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualHubName, id.Name, existing)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update to %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r RouteMapResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RouteMapsClient

			id, err := parse.RouteMapID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := RouteMapModel{
				Name:         id.Name,
				VirtualHubId: parse.NewVirtualHubID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName).ID(),
			}

			if props := resp.RouteMapProperties; props != nil {
				state.Rules = flattenRules(props.Rules)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RouteMapResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.RouteMapsClient

			id, err := parse.RouteMapID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandRules(input []Rule) *[]network.RouteMapRule {
	var rules []network.RouteMapRule
	if input == nil {
		return nil
	}

	for _, v := range input {
		rule := network.RouteMapRule{
			Name:          utils.String(v.Name),
			Actions:       expandActions(v.Actions),
			MatchCriteria: expandCriteria(v.MatchCriteria),
		}

		if v.NextStepIfMatched != "" {
			rule.NextStepIfMatched = v.NextStepIfMatched
		}

		rules = append(rules, rule)
	}

	return &rules
}

func expandActions(input []Action) *[]network.Action {
	var actions []network.Action
	if input == nil {
		return nil
	}

	for _, v := range input {
		action := network.Action{
			Type:       v.Type,
			Parameters: expandParameters(v.Parameters),
		}

		actions = append(actions, action)
	}

	return &actions
}

func expandParameters(input []Parameter) *[]network.Parameter {
	var parameters []network.Parameter
	if input == nil {
		return nil
	}

	for _, item := range input {
		v := item
		parameter := network.Parameter{}

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

func expandCriteria(input []Criterion) *[]network.Criterion {
	var criteria []network.Criterion
	if input == nil {
		return nil
	}

	for _, item := range input {
		v := item
		criterion := network.Criterion{
			MatchCondition: v.MatchCondition,
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

func flattenRules(input *[]network.RouteMapRule) []Rule {
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

		if v.NextStepIfMatched != "" {
			rule.NextStepIfMatched = v.NextStepIfMatched
		}

		rules = append(rules, rule)
	}

	return rules
}

func flattenActions(input *[]network.Action) []Action {
	var actions []Action
	if input == nil {
		return actions
	}

	for _, v := range *input {
		action := Action{
			Parameters: flattenParameters(v.Parameters),
		}

		if v.Type != "" {
			action.Type = v.Type
		}

		actions = append(actions, action)
	}

	return actions
}

func flattenParameters(input *[]network.Parameter) []Parameter {
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

func flattenCriteria(input *[]network.Criterion) []Criterion {
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

		if v.MatchCondition != "" {
			criterion.MatchCondition = v.MatchCondition
		}

		if v.RoutePrefix != nil {
			criterion.RoutePrefix = *v.RoutePrefix
		}

		criteria = append(criteria, criterion)
	}

	return criteria
}
