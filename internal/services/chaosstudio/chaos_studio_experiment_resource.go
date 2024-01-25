package chaosstudio

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ChaosStudioExperimentResource{}
var _ sdk.ResourceWithUpdate = ChaosStudioExperimentResource{}

type ChaosStudioExperimentResource struct{}

func (r ChaosStudioExperimentResource) ModelObject() interface{} {
	return &ChaosStudioExperimentResourceSchema{}
}

type ChaosStudioExperimentResourceSchema struct {
	Identity          []identity.ModelSystemAssignedUserAssigned    `tfschema:"identity"`
	Location          string                                        `tfschema:"location"`
	Name              string                                        `tfschema:"name"`
	ResourceGroupName string                                        `tfschema:"resource_group_name"`
	Selector          []ChaosStudioExperimentResourceSelectorSchema `tfschema:"selector"`
	Step              []ChaosStudioExperimentResourceStepSchema     `tfschema:"step"`
	Tags              map[string]interface{}                        `tfschema:"tags"`
}

type ChaosStudioExperimentResourceSelectorSchema struct {
	//Filter []ChaosStudioExperimentResourceFilterSchema `tfschema:"filter"`
	Id   string `tfschema:"id"`
	Type string `tfschema:"type"`
}

type ChaosStudioExperimentResourceStepSchema struct {
	Branch []ChaosStudioExperimentResourceBranchSchema `tfschema:"branch"`
	Name   string                                      `tfschema:"name"`
}

type ChaosStudioExperimentResourceBranchSchema struct {
	Action []ChaosStudioExperimentResourceActionSchema `tfschema:"action"`
	Name   string                                      `tfschema:"name"`
}

type ChaosStudioExperimentResourceActionSchema struct {
	Name string `tfschema:"name"`
	Type string `tfschema:"type"`
}

type ChaosStudioExperimentResourceFilterSchema struct {
	Type string `tfschema:"type"`
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
		"selector": {
			Required: true,
			Type:     pluginsdk.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Required: true,
						Type:     pluginsdk.TypeString,
					},
					"type": {
						Required: true,
						Type:     pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							"List",
							"Query",
						}, false),
					},
					"filter": {
						MaxItems: 1,
						Optional: true,
						Type:     pluginsdk.TypeList,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Required: true,
									Type:     pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										"Simple",
									}, false),
								},
							},
						},
					},
				},
			},
		},
		"step": {
			Required: true,
			Type:     pluginsdk.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Required: true,
						Type:     pluginsdk.TypeString,
					},
					"branch": {
						Required: true,
						Type:     pluginsdk.TypeList,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Required: true,
									Type:     pluginsdk.TypeString,
								},
								// split into three types of actions?
								"action": {
									Required: true,
									Type:     pluginsdk.TypeList,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											// urn?
											"name": {
												Required: true,
												Type:     pluginsdk.TypeString,
											},
											"type": {
												Required: true,
												Type:     pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice([]string{
													"continuous",
													"delay",
													"discrete",
												}, false),
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
		"tags":     commonschema.Tags(),
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
			payload.Tags = tags.Expand(config.Tags)

			var experimentProperties experiments.ExperimentProperties
			//selector, err := expandSelector(config.Selector)

			if err := r.mapChaosStudioExperimentResourceSchemaToExperimentProperties(config, &experimentProperties); err != nil {
				return fmt.Errorf("flattening steps")
			}

			payload.Properties = experimentProperties
			//if err := r.mapChaosStudioExperimentResourceSchemaToExperiment(config, &payload); err != nil {
			//	return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			//}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

//func expandSelector(input []ChaosStudioExperimentResourceSelectorSchema) ([]experiments.Selector, error) {
//	output := make([]experiments.Selector, 0)
//
//	for _, v := range input {
//		output = append(output, experiments.Selector{
//			Filter: nil,
//			Id:     v.Id,
//			Type:   experiments.SelectorType(v.Type),
//		})
//	}
//
//	return output, nil
//}

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
				if err := r.mapExperimentToChaosStudioExperimentResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
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

			if err := r.mapChaosStudioExperimentResourceSchemaToExperiment(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ChaosStudioExperimentResource) mapChaosStudioExperimentResourceActionSchemaToAction(input ChaosStudioExperimentResourceActionSchema, output *experiments.Action) error {
	output.Name = input.Name
	output.Type = input.Type
	return nil
}

func (r ChaosStudioExperimentResource) mapActionToChaosStudioExperimentResourceActionSchema(input experiments.Action, output *ChaosStudioExperimentResourceActionSchema) error {
	output.Name = input.Name
	output.Type = input.Type
	return nil
}

func (r ChaosStudioExperimentResource) mapChaosStudioExperimentResourceBranchSchemaToBranch(input ChaosStudioExperimentResourceBranchSchema, output *experiments.Branch) error {

	actions := make([]experiments.Action, 0)
	for i, v := range input.Action {
		item := experiments.Action{}
		if err := r.mapChaosStudioExperimentResourceActionSchemaToAction(v, &item); err != nil {
			return fmt.Errorf("mapping ChaosStudioExperimentResourceActionSchema item %d to Action: %+v", i, err)
		}
		actions = append(actions, item)
	}
	output.Actions = actions

	output.Name = input.Name
	return nil
}

func (r ChaosStudioExperimentResource) mapBranchToChaosStudioExperimentResourceBranchSchema(input experiments.Branch, output *ChaosStudioExperimentResourceBranchSchema) error {

	actions := make([]ChaosStudioExperimentResourceActionSchema, 0)
	for i, v := range input.Actions {
		item := ChaosStudioExperimentResourceActionSchema{}
		if err := r.mapActionToChaosStudioExperimentResourceActionSchema(v, &item); err != nil {
			return fmt.Errorf("mapping ChaosStudioExperimentResourceActionSchema item %d to Action: %+v", i, err)
		}
		actions = append(actions, item)
	}
	output.Action = actions

	output.Name = input.Name
	return nil
}

func (r ChaosStudioExperimentResource) mapChaosStudioExperimentResourceFilterSchemaToFilter(input ChaosStudioExperimentResourceFilterSchema, output *experiments.Filter) error {
	output.Type = experiments.FilterType(input.Type)
	return nil
}

func (r ChaosStudioExperimentResource) mapFilterToChaosStudioExperimentResourceFilterSchema(input experiments.Filter, output *ChaosStudioExperimentResourceFilterSchema) error {
	output.Type = string(input.Type)
	return nil
}

func (r ChaosStudioExperimentResource) mapChaosStudioExperimentResourceSchemaToExperiment(input ChaosStudioExperimentResourceSchema, output *experiments.Experiment) error {

	identity, err := identity.ExpandSystemOrUserAssignedMapFromModel(input.Identity)
	if err != nil {
		return fmt.Errorf("expanding SystemOrUserAssigned Identity: %+v", err)
	}
	output.Identity = identity

	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if err := r.mapChaosStudioExperimentResourceSchemaToExperimentProperties(input, &output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "ExperimentProperties", "Properties", err)
	}

	return nil
}

func (r ChaosStudioExperimentResource) mapExperimentToChaosStudioExperimentResourceSchema(input experiments.Experiment, output *ChaosStudioExperimentResourceSchema) error {

	flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMapToModel(input.Identity)
	if err != nil {
		return fmt.Errorf("flattening SystemOrUserAssigned Identity: %+v", err)
	}
	output.Identity = *flattenedIdentity

	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if err := r.mapExperimentPropertiesToChaosStudioExperimentResourceSchema(input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "ExperimentProperties", "Properties", err)
	}

	return nil
}

func (r ChaosStudioExperimentResource) mapChaosStudioExperimentResourceSchemaToExperimentProperties(input ChaosStudioExperimentResourceSchema, output *experiments.ExperimentProperties) error {

	selectors := make([]experiments.Selector, 0)
	for i, v := range input.Selector {
		item := experiments.Selector{}
		if err := r.mapChaosStudioExperimentResourceSelectorSchemaToSelector(v, &item); err != nil {
			return fmt.Errorf("mapping ChaosStudioExperimentResourceSelectorSchema item %d to Selector: %+v", i, err)
		}
		selectors = append(selectors, item)
	}
	output.Selectors = selectors

	steps := make([]experiments.Step, 0)
	for i, v := range input.Step {
		item := experiments.Step{}
		if err := r.mapChaosStudioExperimentResourceStepSchemaToStep(v, &item); err != nil {
			return fmt.Errorf("mapping ChaosStudioExperimentResourceStepSchema item %d to Step: %+v", i, err)
		}
		steps = append(steps, item)
	}
	output.Steps = steps

	return nil
}

func (r ChaosStudioExperimentResource) mapExperimentPropertiesToChaosStudioExperimentResourceSchema(input experiments.ExperimentProperties, output *ChaosStudioExperimentResourceSchema) error {

	selectors := make([]ChaosStudioExperimentResourceSelectorSchema, 0)
	for i, v := range input.Selectors {
		item := ChaosStudioExperimentResourceSelectorSchema{}
		if err := r.mapSelectorToChaosStudioExperimentResourceSelectorSchema(v, &item); err != nil {
			return fmt.Errorf("mapping ChaosStudioExperimentResourceSelectorSchema item %d to Selector: %+v", i, err)
		}
		selectors = append(selectors, item)
	}
	output.Selector = selectors

	steps := make([]ChaosStudioExperimentResourceStepSchema, 0)
	for i, v := range input.Steps {
		item := ChaosStudioExperimentResourceStepSchema{}
		if err := r.mapStepToChaosStudioExperimentResourceStepSchema(v, &item); err != nil {
			return fmt.Errorf("mapping ChaosStudioExperimentResourceStepSchema item %d to Step: %+v", i, err)
		}
		steps = append(steps, item)
	}
	output.Step = steps

	return nil
}

func (r ChaosStudioExperimentResource) mapChaosStudioExperimentResourceSelectorSchemaToSelector(input ChaosStudioExperimentResourceSelectorSchema, output *experiments.Selector) error {
	//if len(input.Filter) > 0 {
	//	if err := r.mapChaosStudioExperimentResourceFilterSchemaToSelector(input.Filter[0], output); err != nil {
	//		return err
	//	}
	//}
	output.Id = input.Id
	output.Type = experiments.SelectorType(input.Type)
	return nil
}

func (r ChaosStudioExperimentResource) mapSelectorToChaosStudioExperimentResourceSelectorSchema(input experiments.Selector, output *ChaosStudioExperimentResourceSelectorSchema) error {
	//tmpFilter := &ChaosStudioExperimentResourceFilterSchema{}
	//if err := r.mapSelectorToChaosStudioExperimentResourceFilterSchema(input, tmpFilter); err != nil {
	//	return err
	//} else {
	//	output.Filter = make([]ChaosStudioExperimentResourceFilterSchema, 0)
	//	output.Filter = append(output.Filter, *tmpFilter)
	//}
	output.Id = input.Id
	output.Type = string(input.Type)
	return nil
}

func (r ChaosStudioExperimentResource) mapChaosStudioExperimentResourceStepSchemaToStep(input ChaosStudioExperimentResourceStepSchema, output *experiments.Step) error {

	branches := make([]experiments.Branch, 0)
	for i, v := range input.Branch {
		item := experiments.Branch{}
		if err := r.mapChaosStudioExperimentResourceBranchSchemaToBranch(v, &item); err != nil {
			return fmt.Errorf("mapping ChaosStudioExperimentResourceBranchSchema item %d to Branch: %+v", i, err)
		}
		branches = append(branches, item)
	}
	output.Branches = branches

	output.Name = input.Name
	return nil
}

func (r ChaosStudioExperimentResource) mapStepToChaosStudioExperimentResourceStepSchema(input experiments.Step, output *ChaosStudioExperimentResourceStepSchema) error {

	branches := make([]ChaosStudioExperimentResourceBranchSchema, 0)
	for i, v := range input.Branches {
		item := ChaosStudioExperimentResourceBranchSchema{}
		if err := r.mapBranchToChaosStudioExperimentResourceBranchSchema(v, &item); err != nil {
			return fmt.Errorf("mapping ChaosStudioExperimentResourceBranchSchema item %d to Branch: %+v", i, err)
		}
		branches = append(branches, item)
	}
	output.Branch = branches

	output.Name = input.Name
	return nil
}
