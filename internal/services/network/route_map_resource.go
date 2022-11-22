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
	"github.com/tombuildsstuff/kermit/sdk/network/2022-05-01/network"
)

type RouteMapModel struct {
	Name                          string   `tfschema:"name"`
	VirtualHubId                  string   `tfschema:"virtual_hub_id"`
	AssociatedInboundConnections  []string `tfschema:"associated_inbound_connections"`
	AssociatedOutboundConnections []string `tfschema:"associated_outbound_connections"`
	Rules                         []Rule   `tfschema:"rules"`
}

type Rule struct {
	Actions           []Action         `tfschema:"actions"`
	MatchCriteria     []Criterion      `tfschema:"match_criteria"`
	Name              string           `tfschema:"name"`
	NextStepIfMatched network.NextStep `tfschema:"next_step_if_matched"`
}

type Action struct {
	Parameters []Parameter                `tfschema:"parameters"`
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
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"virtual_hub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualHubID,
		},

		"associated_inbound_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"associated_outbound_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"rules": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"parameters": {
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
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(network.RouteMapActionTypeReplace),
										string(network.RouteMapActionTypeDrop),
										string(network.RouteMapActionTypeUnknown),
										string(network.RouteMapActionTypeRemove),
										string(network.RouteMapActionTypeAdd),
									}, false),
								},
							},
						},
					},

					"match_criteria": {
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

								"match_condition": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(network.RouteMapMatchConditionEquals),
										string(network.RouteMapMatchConditionNotContains),
										string(network.RouteMapMatchConditionNotEquals),
										string(network.RouteMapMatchConditionUnknown),
										string(network.RouteMapMatchConditionContains),
									}, false),
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

					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"next_step_if_matched": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(network.NextStepUnknown),
							string(network.NextStepContinue),
							string(network.NextStepTerminate),
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
					AssociatedInboundConnections:  &model.AssociatedInboundConnections,
					AssociatedOutboundConnections: &model.AssociatedOutboundConnections,
				},
			}

			rules, err := expandRules(model.Rules)
			if err != nil {
				return err
			}
			props.RouteMapProperties.Rules = rules

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

			if metadata.ResourceData.HasChange("associated_inbound_connections") {
				existing.RouteMapProperties.AssociatedInboundConnections = &model.AssociatedInboundConnections
			}

			if metadata.ResourceData.HasChange("associated_outbound_connections") {
				existing.RouteMapProperties.AssociatedOutboundConnections = &model.AssociatedOutboundConnections
			}

			if metadata.ResourceData.HasChange("rules") {
				rules, err := expandRules(model.Rules)
				if err != nil {
					return err
				}
				existing.RouteMapProperties.Rules = rules
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
				if props.AssociatedInboundConnections != nil {
					state.AssociatedInboundConnections = *props.AssociatedInboundConnections
				}

				if props.AssociatedOutboundConnections != nil {
					state.AssociatedOutboundConnections = *props.AssociatedOutboundConnections
				}

				rules, err := flattenRules(props.Rules)
				if err != nil {
					return err
				}
				state.Rules = rules
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

func expandRules(input []Rule) (*[]network.RouteMapRule, error) {
	var rules []network.RouteMapRule

	for _, v := range input {
		rule := network.RouteMapRule{
			NextStepIfMatched: v.NextStepIfMatched,
		}

		actions, err := expandActions(v.Actions)
		if err != nil {
			return nil, err
		}
		rule.Actions = actions

		matchCriteria, err := expandCriterions(v.MatchCriteria)
		if err != nil {
			return nil, err
		}
		rule.MatchCriteria = matchCriteria

		if v.Name != "" {
			rule.Name = &v.Name
		}

		rules = append(rules, rule)
	}

	return &rules, nil
}

func expandActions(input []Action) (*[]network.Action, error) {
	var actions []network.Action

	for _, v := range input {
		action := network.Action{
			Type: v.Type,
		}

		parameters, err := expandParameters(v.Parameters)
		if err != nil {
			return nil, err
		}
		action.Parameters = parameters

		actions = append(actions, action)
	}

	return &actions, nil
}

func expandParameters(input []Parameter) (*[]network.Parameter, error) {
	var paramters []network.Parameter

	for _, v := range input {
		parameter := network.Parameter{
			AsPath:      &v.AsPath,
			Community:   &v.Community,
			RoutePrefix: &v.RoutePrefix,
		}

		paramters = append(paramters, parameter)
	}

	return &paramters, nil
}

func expandCriterions(input []Criterion) (*[]network.Criterion, error) {
	var criterions []network.Criterion

	for _, v := range input {
		criterion := network.Criterion{
			AsPath:         &v.AsPath,
			Community:      &v.Community,
			MatchCondition: v.MatchCondition,
			RoutePrefix:    &v.RoutePrefix,
		}

		criterions = append(criterions, criterion)
	}

	return &criterions, nil
}

func flattenRules(input *[]network.RouteMapRule) ([]Rule, error) {
	var rules []Rule
	if input == nil {
		return rules, nil
	}

	for _, v := range *input {
		rule := Rule{}

		actions, err := flattenActions(v.Actions)
		if err != nil {
			return nil, err
		}
		rule.Actions = actions

		matchCriteria, err := flattenCriterions(v.MatchCriteria)
		if err != nil {
			return nil, err
		}
		rule.MatchCriteria = matchCriteria

		if v.Name != nil {
			rule.Name = *v.Name
		}

		if v.NextStepIfMatched != "" {
			rule.NextStepIfMatched = v.NextStepIfMatched
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func flattenActions(input *[]network.Action) ([]Action, error) {
	var actions []Action
	if input == nil {
		return actions, nil
	}

	for _, v := range *input {
		action := Action{}

		parameters, err := flattenParameters(v.Parameters)
		if err != nil {
			return nil, err
		}
		action.Parameters = parameters

		if v.Type != "" {
			action.Type = v.Type
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func flattenParameters(input *[]network.Parameter) ([]Parameter, error) {
	var parameters []Parameter
	if input == nil {
		return parameters, nil
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

	return parameters, nil
}

func flattenCriterions(input *[]network.Criterion) ([]Criterion, error) {
	var criterions []Criterion
	if input == nil {
		return criterions, nil
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

		criterions = append(criterions, criterion)
	}

	return criterions, nil
}
