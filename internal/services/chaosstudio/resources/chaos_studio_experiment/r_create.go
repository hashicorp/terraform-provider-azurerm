// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_experiment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
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
