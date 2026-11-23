// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_experiment

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments"
)

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
