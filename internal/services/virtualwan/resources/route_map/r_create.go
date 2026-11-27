// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_map

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RouteMapModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := virtualwanClient.From(metadata.Client).VirtualWANs
			virtualHubId, err := virtualwans.ParseVirtualHubID(model.VirtualHubId)
			if err != nil {
				return err
			}

			id := virtualwans.NewRouteMapID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroupName, virtualHubId.VirtualHubName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.RouteMapsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			props := &virtualwans.RouteMap{
				Properties: &virtualwans.RouteMapProperties{
					Rules: expandRules(model.Rules),
				},
			}

			if err := client.RouteMapsCreateOrUpdateCallbackThenPoll(ctx, id, *props, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func expandRules(input []Rule) *[]virtualwans.RouteMapRule {
	if input == nil {
		return nil
	}

	rules := make([]virtualwans.RouteMapRule, 0, len(input))
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
	if input == nil {
		return nil
	}

	actions := make([]virtualwans.Action, 0, len(input))
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
	if input == nil {
		return nil
	}

	parameters := make([]virtualwans.Parameter, 0, len(input))
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
	if input == nil {
		return nil
	}

	criteria := make([]virtualwans.Criterion, 0, len(input))
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
