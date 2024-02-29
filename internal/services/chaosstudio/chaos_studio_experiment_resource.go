package chaosstudio

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/chaosstudio/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ChaosStudioExperimentResource{}
var _ sdk.ResourceWithUpdate = ChaosStudioExperimentResource{}

const (
	continuousActionType = "continuous"
	delayActionType      = "delay"
	discreteActionType   = "discrete"
)

type ChaosStudioExperimentResource struct{}

func (r ChaosStudioExperimentResource) ModelObject() interface{} {
	return &ChaosStudioExperimentResourceSchema{}
}

type ChaosStudioExperimentResourceSchema struct {
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location          string                                     `tfschema:"location"`
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Selectors         []SelectorSchema                           `tfschema:"selectors"`
	Steps             []StepSchema                               `tfschema:"steps"`
	// tags are not fully supported yet, you can send them to the API, but they won't be returned
	// Tags              map[string]interface{}                     `tfschema:"tags"`
}

type SelectorSchema struct {
	Name      string   `tfschema:"name"`
	TargetIds []string `tfschema:"chaos_studio_target_ids"`
}

type StepSchema struct {
	Branch []BranchSchema `tfschema:"branch"`
	Name   string         `tfschema:"name"`
}

type BranchSchema struct {
	Actions []ActionSchema `tfschema:"actions"`
	Name    string         `tfschema:"name"`
}

type ActionSchema struct {
	ActionType   string            `tfschema:"action_type"`
	SelectorName string            `tfschema:"selector_name"`
	Duration     string            `tfschema:"duration"`
	Urn          string            `tfschema:"urn"`
	Parameters   map[string]string `tfschema:"parameters"`
}

func (r ChaosStudioExperimentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return experiments.ValidateExperimentID
}

func (r ChaosStudioExperimentResource) ResourceType() string {
	return "azurerm_chaos_studio_experiment"
}

func (r ChaosStudioExperimentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"selectors": {
			Required: true,
			Type:     pluginsdk.TypeList,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Required:     true,
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"chaos_studio_target_ids": {
						Required: true,
						Type:     pluginsdk.TypeList,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: commonids.ValidateChaosStudioTargetID,
						},
					},
				},
			},
		},
		"steps": {
			Required: true,
			Type:     pluginsdk.TypeList,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Required:     true,
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"branch": {
						Required: true,
						Type:     pluginsdk.TypeList,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Required:     true,
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"actions": {
									Required: true,
									Type:     pluginsdk.TypeList,
									MinItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"action_type": {
												Required: true,
												Type:     pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice([]string{
													continuousActionType,
													delayActionType,
													discreteActionType,
												}, false),
											},
											// the different types of actions require different properties to be set
											// which is why the validation for these is done in expandActions
											"urn": {
												Optional: true,
												Type:     pluginsdk.TypeString,
											},
											"selector_name": {
												Optional: true,
												Type:     pluginsdk.TypeString,
											},
											"duration": {
												Optional: true,
												Type:     pluginsdk.TypeString,
											},
											"parameters": {
												Type:     pluginsdk.TypeMap,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),
	}
}

func (r ChaosStudioExperimentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ChaosStudioExperimentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Experiments

			var config ChaosStudioExperimentResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := experiments.NewExperimentID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload experiments.Experiment

			expandedIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding SystemOrUserAssigned Identity: %+v", err)
			}
			payload.Identity = expandedIdentity

			payload.Location = location.Normalize(config.Location)

			var experimentProperties experiments.ExperimentProperties

			selectors, err := expandSelectors(config.Selectors)
			if err != nil {
				return fmt.Errorf("expanding `selectors`: %+v", err)
			}
			experimentProperties.Selectors = *selectors

			steps, err := expandSteps(config.Steps)
			if err != nil {
				return fmt.Errorf("expanding `steps`: %+v", err)
			}
			experimentProperties.Steps = *steps

			payload.Properties = experimentProperties

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ChaosStudioExperimentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Experiments
			schema := ChaosStudioExperimentResourceSchema{}

			id, err := experiments.ParseExperimentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.ExperimentName
				schema.ResourceGroupName = id.ResourceGroupName
				schema.Location = location.Normalize(model.Location)

				props := model.Properties

				selectors, err := flattenSelector(props.Selectors)
				if err != nil {
					return fmt.Errorf("flattening `selectors`: %+v", err)
				}
				schema.Selectors = pointer.From(selectors)

				steps, err := flattenSteps(props.Steps)
				if err != nil {
					return fmt.Errorf("flattening `steps`: %+v", err)
				}
				schema.Steps = pointer.From(steps)

				flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				schema.Identity = *flattenedIdentity
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ChaosStudioExperimentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Experiments

			id, err := experiments.ParseExperimentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ChaosStudioExperimentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Experiments

			id, err := experiments.ParseExperimentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ChaosStudioExperimentResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemOrUserAssignedMapFromModel(config.Identity)
				if err != nil {
					return fmt.Errorf("expanding SystemOrUserAssigned Identity: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("selectors") {
				selectors, err := expandSelectors(config.Selectors)
				if err != nil {
					return fmt.Errorf("expanding `selectors`: %+v", err)
				}
				payload.Properties.Selectors = *selectors
			}

			if metadata.ResourceData.HasChange("steps") {
				steps, err := expandSteps(config.Steps)
				if err != nil {
					return fmt.Errorf("expanding `steps`: %+v", err)
				}
				payload.Properties.Steps = *steps
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			// the PUT method for updates returns a 200 instead of a 201/202 which means we don't build the poller property
			// this can be removed when https://github.com/Azure/azure-rest-api-specs/issues/27659 is fixed
			pollerType := custompollers.NewChaosStudioExperimentPoller(client, *id)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandSelectors(input []SelectorSchema) (*[]experiments.Selector, error) {
	output := make([]experiments.Selector, 0)

	for _, v := range input {
		targetsOutput := make([]experiments.TargetReference, 0)
		for _, t := range v.TargetIds {
			targetId, err := commonids.ParseChaosStudioTargetID(t)
			if err != nil {
				return nil, err
			}
			targetsOutput = append(targetsOutput, experiments.TargetReference{
				Id:   targetId.ID(),
				Type: experiments.TargetReferenceTypeChaosTarget,
			})
		}
		output = append(output, experiments.ListSelector{
			Targets: targetsOutput,
			Filter:  nil,
			Id:      v.Name,
		})
	}
	return &output, nil
}

func expandSteps(input []StepSchema) (*[]experiments.Step, error) {
	output := make([]experiments.Step, 0)

	for _, step := range input {
		branches := make([]experiments.Branch, 0)
		for _, branch := range step.Branch {
			actions, err := expandActions(branch.Actions)
			if err != nil {
				return nil, fmt.Errorf("expanding `actions`: %+v", err)
			}
			branches = append(branches, experiments.Branch{
				Actions: *actions,
				Name:    branch.Name,
			})
		}
		output = append(output, experiments.Step{
			Name:     step.Name,
			Branches: branches,
		})
	}

	return &output, nil
}

func expandActions(input []ActionSchema) (*[]experiments.Action, error) {
	output := make([]experiments.Action, 0)

	for _, action := range input {
		parameters := make([]experiments.KeyValuePair, 0)
		if len(action.Parameters) > 0 {
			for k, v := range action.Parameters {
				parameters = append(parameters, experiments.KeyValuePair{
					Key:   k,
					Value: v,
				})
			}
		}

		switch action.ActionType {
		case continuousActionType:
			if action.Duration == "" || action.SelectorName == "" || action.Urn == "" {
				return nil, fmt.Errorf("`duration`, `selector_name` and `urn` must be set for actions with `action_type` of `continuous`")
			}
			output = append(output, experiments.ContinuousAction{
				Duration:   action.Duration,
				Parameters: parameters,
				SelectorId: action.SelectorName,
				Name:       action.Urn,
			})
		case delayActionType:
			if action.Duration == "" {
				return nil, fmt.Errorf("`duration` must be set for actions with `action_type` of `delay`")
			}
			output = append(output, experiments.DelayAction{
				Duration: action.Duration,
				Name:     "urn:csci:microsoft:chaosStudio:timedDelay/1.0",
			})
		case discreteActionType:
			if action.SelectorName == "" || action.Urn == "" {
				return nil, fmt.Errorf("`selector_name` and `urn` must be set for actions with `action_type` of `discrete`")
			}
			output = append(output, experiments.DiscreteAction{
				Parameters: parameters,
				SelectorId: action.SelectorName,
				Name:       action.Urn,
			})
		}
	}

	return &output, nil
}

func flattenSelector(input []experiments.Selector) (*[]SelectorSchema, error) {
	output := make([]SelectorSchema, 0)

	if len(input) == 0 {
		return &output, nil
	}

	for _, selector := range input {
		targetIds := make([]string, 0)
		ls, ok := selector.(experiments.ListSelector)
		if !ok {
			return nil, fmt.Errorf("selector is not of type ListSelector")
		}
		for _, t := range ls.Targets {
			targetIds = append(targetIds, t.Id)
		}
		output = append(output, SelectorSchema{
			Name:      ls.Id,
			TargetIds: targetIds,
		})
	}

	return &output, nil
}

func flattenSteps(input []experiments.Step) (*[]StepSchema, error) {
	output := make([]StepSchema, 0)

	if len(input) == 0 {
		return &output, nil
	}

	for _, step := range input {
		branches := make([]BranchSchema, 0)
		for _, branch := range step.Branches {
			actions, err := flattenActions(branch.Actions)
			if err != nil {
				return nil, fmt.Errorf("flattening `actions`: %+v", err)
			}
			branches = append(branches, BranchSchema{
				Actions: *actions,
				Name:    branch.Name,
			})
		}
		output = append(output, StepSchema{
			Branch: branches,
			Name:   step.Name,
		})
	}

	return &output, nil
}

func flattenActions(input []experiments.Action) (*[]ActionSchema, error) {
	output := make([]ActionSchema, 0)

	if len(input) == 0 {
		return &output, nil
	}

	for _, action := range input {
		actionOutput := ActionSchema{}

		switch a := action.(type) {
		case experiments.ContinuousAction:
			parameters := make(map[string]string)
			for _, p := range a.Parameters {
				parameters[p.Key] = p.Value
			}
			actionOutput.Parameters = parameters
			actionOutput.SelectorName = a.SelectorId
			actionOutput.Urn = a.Name
			actionOutput.Duration = a.Duration
			actionOutput.ActionType = continuousActionType
		case experiments.DelayAction:
			actionOutput.Duration = a.Duration
			actionOutput.ActionType = delayActionType
		case experiments.DiscreteAction:
			parameters := make(map[string]string)
			for _, p := range a.Parameters {
				parameters[p.Key] = p.Value
			}
			actionOutput.Parameters = parameters
			actionOutput.SelectorName = a.SelectorId
			actionOutput.Urn = a.Name
			actionOutput.ActionType = discreteActionType
		default:
			return nil, fmt.Errorf("action is not of type `continuous`, `delay` or `discrete`")
		}

		output = append(output, actionOutput)
	}

	return &output, nil
}
