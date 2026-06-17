// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	legacyrulesets "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	batchRules "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
)

var _ pollers.PollerType = &frontDoorBatchRuleSetUpdatePoller{}

const frontDoorBatchRuleSetPollInterval = 30 * time.Second

type frontDoorBatchRuleSetUpdatePoller struct {
	client          *azuresdkhacks.BatchRuleSetsClient
	id              legacyrulesets.RuleSetId
	input           azuresdkhacks.BatchRuleSetResource
	operationIssued bool
}

func NewFrontDoorBatchRuleSetUpdatePoller(client *azuresdkhacks.BatchRuleSetsClient, id legacyrulesets.RuleSetId, input azuresdkhacks.BatchRuleSetResource) pollers.PollerType {
	return &frontDoorBatchRuleSetUpdatePoller{
		client: client,
		id:     id,
		input:  input,
	}
}

func (p *frontDoorBatchRuleSetUpdatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if !p.operationIssued {
		result, err := p.client.Create(ctx, p.id, p.input)
		if err != nil {
			return nil, fmt.Errorf("updating %s: %+v", p.id.ID(), err)
		}

		if err := result.Poller.PollUntilDone(ctx); err != nil {
			return nil, fmt.Errorf("waiting for the update of %s: %+v", p.id.ID(), err)
		}

		p.operationIssued = true
	}

	ready, err := frontDoorBatchRuleSetSettledForUpdate(ctx, p.client, p.id, p.input)
	if err != nil {
		return nil, err
	}
	if !ready {
		return frontDoorBatchRuleSetInProgress(), nil
	}

	return frontDoorBatchRuleSetSucceeded(), nil
}

func frontDoorBatchRuleSetSettledForUpdate(ctx context.Context, client *azuresdkhacks.BatchRuleSetsClient, id legacyrulesets.RuleSetId, desired azuresdkhacks.BatchRuleSetResource) (bool, error) {
	resp, err := client.Get(ctx, id)
	if err != nil {
		return false, fmt.Errorf("retrieving %s while waiting for batch rule set state to settle after update: %+v", id.ID(), err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return false, nil
	}

	ready, err := batchRuleSetStatusesSettled(resp.Model)
	if err != nil {
		return false, fmt.Errorf("waiting for %s: %+v", id.ID(), err)
	}
	if !ready {
		return false, nil
	}

	matches, err := batchRuleSetOriginGroupOverridesMatchDesired(resp.Model, desired)
	if err != nil {
		return false, fmt.Errorf("checking origin-group dissociation for %s: %+v", id.ID(), err)
	}
	if !matches {
		return false, nil
	}

	return true, nil
}

func batchRuleSetStatusesSettled(input *azuresdkhacks.BatchRuleSetResource) (bool, error) {
	if input == nil || input.Properties == nil {
		return false, nil
	}

	properties := input.Properties
	deploymentStatus := batchRuleSetStringValue(properties.DeploymentStatus)
	provisioningState := batchRuleSetStringValue(properties.ProvisioningState)

	if strings.EqualFold(provisioningState, string(batchRules.AfdProvisioningStateFailed)) {
		return false, fmt.Errorf("batch rule set entered failed state with `deploymentStatus` `%s` and `provisioningState` `%s`", deploymentStatus, provisioningState)
	}

	// For batch rulesets, `provisioningState` is the authoritative top-level
	// readiness signal. `deploymentStatus` is still useful for diagnostics, but it
	// can legitimately remain `NotStarted` even after the service has applied the
	// update.
	if !strings.EqualFold(provisioningState, string(batchRules.AfdProvisioningStateSucceeded)) {
		return false, nil
	}

	if properties.Rules == nil {
		return true, nil
	}

	for _, rule := range *properties.Rules {
		ruleName, err := batchRuleSetRuleName(rule)
		if err != nil {
			return false, err
		}

		ruleDeploymentStatus := batchRuleDeploymentStatusValue(rule.DeploymentStatus)
		ruleProvisioningState := batchRuleProvisioningStateValue(rule.ProvisioningState)

		if strings.EqualFold(ruleProvisioningState, string(batchRules.AfdProvisioningStateFailed)) {
			return false, fmt.Errorf("rule `%s` entered failed state with `deploymentStatus` `%s` and `provisioningState` `%s`", ruleName, ruleDeploymentStatus, ruleProvisioningState)
		}

		if ruleProvisioningState != "" && !strings.EqualFold(ruleProvisioningState, string(batchRules.AfdProvisioningStateSucceeded)) {
			return false, nil
		}
	}

	return true, nil
}

func batchRuleSetOriginGroupOverridesMatchDesired(actual *azuresdkhacks.BatchRuleSetResource, desired azuresdkhacks.BatchRuleSetResource) (bool, error) {
	var actualRules *[]azuresdkhacks.BatchRuleProperties
	if actual != nil && actual.Properties != nil {
		actualRules = actual.Properties.Rules
	}

	actualTargets, err := batchRuleSetRuleOriginGroupTargets(actualRules)
	if err != nil {
		return false, err
	}

	var desiredRules *[]azuresdkhacks.BatchRuleProperties
	if desired.Properties != nil {
		desiredRules = desired.Properties.Rules
	}

	desiredTargets, err := batchRuleSetRuleOriginGroupTargets(desiredRules)
	if err != nil {
		return false, err
	}

	if len(actualTargets) != len(desiredTargets) {
		return false, nil
	}

	for ruleName, desiredTarget := range desiredTargets {
		actualTarget, ok := actualTargets[ruleName]
		if !ok {
			return false, nil
		}

		if !strings.EqualFold(actualTarget, desiredTarget) {
			return false, nil
		}
	}

	return true, nil
}

func batchRuleSetRuleOriginGroupTargets(input *[]azuresdkhacks.BatchRuleProperties) (map[string]string, error) {
	results := make(map[string]string)
	if input == nil {
		return results, nil
	}

	for _, rule := range *input {
		ruleName, err := batchRuleSetRuleName(rule)
		if err != nil {
			return nil, err
		}

		target, err := batchRuleSetRuleOriginGroupTarget(rule.Actions)
		if err != nil {
			return nil, fmt.Errorf("reading `route_configuration_override_action` for rule `%s`: %+v", ruleName, err)
		}

		results[ruleName] = target
	}

	return results, nil
}

func batchRuleSetRuleName(input azuresdkhacks.BatchRuleProperties) (string, error) {
	if input.Name != nil && *input.Name != "" {
		return *input.Name, nil
	}

	if input.RuleName != nil && *input.RuleName != "" {
		return *input.RuleName, nil
	}

	return "", errors.New("expected each batch rule to contain a non-empty `name`")
}

func batchRuleSetRuleOriginGroupTarget(actions *[]batchRules.DeliveryRuleAction) (string, error) {
	if actions == nil {
		return "", nil
	}

	for _, action := range *actions {
		if action.DeliveryRuleAction().Name != batchRules.DeliveryRuleActionNameRouteConfigurationOverride {
			continue
		}

		routeAction, ok := action.(batchRules.DeliveryRuleRouteConfigurationOverrideAction)
		if !ok {
			return "", fmt.Errorf("expected DeliveryRuleRouteConfigurationOverrideAction but got %T", action)
		}

		originGroupOverride := routeAction.Parameters.OriginGroupOverride
		if originGroupOverride == nil || originGroupOverride.OriginGroup == nil || originGroupOverride.OriginGroup.Id == nil {
			return "", nil
		}

		originGroup, err := parse.FrontDoorOriginGroupIDInsensitively(*originGroupOverride.OriginGroup.Id)
		if err != nil {
			return "", fmt.Errorf("parsing `originGroupOverride`: %+v", err)
		}

		return originGroup.ID(), nil
	}

	return "", nil
}

func batchRuleSetStringValue(input *string) string {
	if input == nil {
		return ""
	}

	return *input
}

func batchRuleDeploymentStatusValue(input *batchRules.DeploymentStatus) string {
	if input == nil {
		return ""
	}

	return string(*input)
}

func batchRuleProvisioningStateValue(input *batchRules.AfdProvisioningState) string {
	if input == nil {
		return ""
	}

	return string(*input)
}

func frontDoorBatchRuleSetInProgress() *pollers.PollResult {
	return &pollers.PollResult{
		PollInterval: frontDoorBatchRuleSetPollInterval,
		Status:       pollers.PollingStatusInProgress,
	}
}

func frontDoorBatchRuleSetSucceeded() *pollers.PollResult {
	return &pollers.PollResult{
		PollInterval: frontDoorBatchRuleSetPollInterval,
		Status:       pollers.PollingStatusSucceeded,
	}
}
