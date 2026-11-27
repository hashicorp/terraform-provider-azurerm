// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_experiment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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

const (
	continuousActionType = "continuous"
	delayActionType      = "delay"
	discreteActionType   = "discrete"
)

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
