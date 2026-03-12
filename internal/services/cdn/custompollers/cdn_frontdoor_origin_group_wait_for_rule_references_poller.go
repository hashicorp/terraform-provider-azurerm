// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-09-01/rules"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
)

var _ pollers.PollerType = &frontDoorOriginGroupWaitForRuleReferencesPoller{}

type frontDoorOriginGroupWaitForRuleReferencesPoller struct {
	ruleSetsClient *rulesets.RuleSetsClient
	rulesClient    *rules.RulesClient
	id             parse.FrontDoorOriginGroupId
}

var (
	frontDoorOriginGroupWaitForRuleReferencesSuccess = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	frontDoorOriginGroupWaitForRuleReferencesInProgress = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewFrontDoorOriginGroupWaitForRuleReferencesPoller(ruleSetsClient *rulesets.RuleSetsClient, rulesClient *rules.RulesClient, id parse.FrontDoorOriginGroupId) pollers.PollerType {
	return &frontDoorOriginGroupWaitForRuleReferencesPoller{
		ruleSetsClient: ruleSetsClient,
		rulesClient:    rulesClient,
		id:             id,
	}
}

func (p frontDoorOriginGroupWaitForRuleReferencesPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	profileId := rulesets.NewProfileID(p.id.SubscriptionId, p.id.ResourceGroup, p.id.ProfileName)
	ruleSetsResult, err := p.ruleSetsClient.ListByProfileComplete(ctx, profileId)
	if err != nil {
		if response.WasNotFound(ruleSetsResult.LatestHttpResponse) {
			return &frontDoorOriginGroupWaitForRuleReferencesSuccess, nil
		}

		return nil, fmt.Errorf("listing rule sets for %s: %+v", p.id, err)
	}

	for _, ruleSet := range ruleSetsResult.Items {
		ruleSetName := pointer.From(ruleSet.Name)
		if ruleSetName == "" {
			continue
		}

		ruleSetId := rules.NewRuleSetID(p.id.SubscriptionId, p.id.ResourceGroup, p.id.ProfileName, ruleSetName)
		rulesResult, err := p.rulesClient.ListByRuleSetComplete(ctx, ruleSetId)
		if err != nil {
			if response.WasNotFound(rulesResult.LatestHttpResponse) {
				continue
			}

			return nil, fmt.Errorf("listing rules for %s: %+v", ruleSetId, err)
		}

		for _, rule := range rulesResult.Items {
			if referencesOriginGroup(rule, p.id) {
				ruleName := pointer.From(rule.Name)
				log.Printf("[DEBUG] AFD Origin Group %s still referenced by rule set %q rule %q", p.id, ruleSetName, ruleName)
				return &frontDoorOriginGroupWaitForRuleReferencesInProgress, nil
			}
		}
	}

	return &frontDoorOriginGroupWaitForRuleReferencesSuccess, nil
}

func referencesOriginGroup(rule rules.Rule, originGroupId parse.FrontDoorOriginGroupId) bool {
	if rule.Properties == nil || rule.Properties.Actions == nil {
		return false
	}

	for _, action := range *rule.Properties.Actions {
		switch v := action.(type) {
		case rules.DeliveryRuleRouteConfigurationOverrideAction:
			if routeConfigurationOverrideReferencesOriginGroup(v.Parameters.OriginGroupOverride, originGroupId) {
				return true
			}
		case rules.OriginGroupOverrideAction:
			if v.Parameters.OriginGroup.Id == nil {
				continue
			}

			if originGroupIdsMatch(*v.Parameters.OriginGroup.Id, originGroupId) {
				return true
			}
		}
	}

	return false
}

func routeConfigurationOverrideReferencesOriginGroup(override *rules.OriginGroupOverride, originGroupId parse.FrontDoorOriginGroupId) bool {
	if override == nil || override.OriginGroup == nil || override.OriginGroup.Id == nil {
		return false
	}

	return originGroupIdsMatch(*override.OriginGroup.Id, originGroupId)
}

func originGroupIdsMatch(candidate string, originGroupId parse.FrontDoorOriginGroupId) bool {
	parsedCandidate, err := parse.FrontDoorOriginGroupIDInsensitively(candidate)
	if err != nil {
		return strings.EqualFold(candidate, originGroupId.ID())
	}

	return strings.EqualFold(parsedCandidate.ID(), originGroupId.ID())
}
