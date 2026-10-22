// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_activity_log_alert

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandMonitorActivityLogAlertCriteria(input []any) activitylogalertsapis.AlertRuleAllOfCondition {
	conditions := make([]activitylogalertsapis.AlertRuleAnyOfOrLeafCondition, 0)
	v := input[0].(map[string]any)

	if category := v["category"].(string); category != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("category"),
			Equals: pointer.To(category),
		})
	}

	if op := v["operation_name"].(string); op != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("operationName"),
			Equals: pointer.To(op),
		})
	}

	if caller := v["caller"].(string); caller != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("caller"),
			Equals: pointer.To(caller),
		})
	}

	if level := v["level"].(string); level != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("level"),
			Equals: pointer.To(level),
		})
	}

	if levels := v["levels"].([]any); len(levels) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(levels, "level"),
		})
	}

	if resourceProvider := v["resource_provider"].(string); resourceProvider != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("resourceProvider"),
			Equals: pointer.To(resourceProvider),
		})
	}

	if resourceProviders := v["resource_providers"].([]any); len(resourceProviders) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(resourceProviders, "resourceProvider"),
		})
	}

	if resourceType := v["resource_type"].(string); resourceType != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("resourceType"),
			Equals: pointer.To(resourceType),
		})
	}

	if resourceTypes := v["resource_types"].([]any); len(resourceTypes) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(resourceTypes, "resourceType"),
		})
	}

	if resourceGroup := v["resource_group"].(string); resourceGroup != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("resourceGroup"),
			Equals: pointer.To(resourceGroup),
		})
	}

	if resourceGroups := v["resource_groups"].([]any); len(resourceGroups) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(resourceGroups, "resourceGroup"),
		})
	}

	if id := v["resource_id"].(string); id != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("resourceId"),
			Equals: pointer.To(id),
		})
	}

	if resourceIds := v["resource_ids"].([]any); len(resourceIds) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(resourceIds, "resourceId"),
		})
	}

	if status := v["status"].(string); status != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("status"),
			Equals: pointer.To(status),
		})
	}

	if statuses := v["statuses"].([]any); len(statuses) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(statuses, "status"),
		})
	}

	if subStatus := v["sub_status"].(string); subStatus != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("subStatus"),
			Equals: pointer.To(subStatus),
		})
	}

	if statuses := v["sub_statuses"].([]any); len(statuses) > 0 {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			AnyOf: expandAnyOfCondition(statuses, "subStatus"),
		})
	}

	if recommendationType := v["recommendation_type"].(string); recommendationType != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("properties.recommendationType"),
			Equals: pointer.To(recommendationType),
		})
	}

	if recommendationCategory := v["recommendation_category"].(string); recommendationCategory != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("properties.recommendationCategory"),
			Equals: pointer.To(recommendationCategory),
		})
	}

	if recommendationImpact := v["recommendation_impact"].(string); recommendationImpact != "" {
		conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
			Field:  pointer.To("properties.recommendationImpact"),
			Equals: pointer.To(recommendationImpact),
		})
	}

	if resourceHealth := v["resource_health"].([]any); len(resourceHealth) > 0 {
		conditions = expandResourceHealth(resourceHealth, conditions)
	}

	if serviceHealth := v["service_health"].([]any); len(serviceHealth) > 0 {
		conditions = expandServiceHealth(serviceHealth, conditions)
	}

	return activitylogalertsapis.AlertRuleAllOfCondition{
		AllOf: conditions,
	}
}

func expandAnyOfCondition(input []any, field string) *[]activitylogalertsapis.AlertRuleLeafCondition {
	conditions := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
	for _, v := range input {
		conditions = append(conditions, activitylogalertsapis.AlertRuleLeafCondition{
			Field:  pointer.To(field),
			Equals: pointer.To(v.(string)),
		})
	}
	return &conditions
}

func expandResourceHealth(resourceHealth []any, conditions []activitylogalertsapis.AlertRuleAnyOfOrLeafCondition) []activitylogalertsapis.AlertRuleAnyOfOrLeafCondition {
	for _, serviceItem := range resourceHealth {
		if serviceItem == nil {
			continue
		}
		vs := serviceItem.(map[string]any)

		cv := vs["current"].(*pluginsdk.Set)
		if len(cv.List()) > 0 {
			ruleLeafCondition := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
			for _, e := range cv.List() {
				event := e.(string)
				ruleLeafCondition = append(ruleLeafCondition, activitylogalertsapis.AlertRuleLeafCondition{
					Field:  pointer.To("properties.currentHealthStatus"),
					Equals: pointer.To(event),
				})
			}
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				AnyOf: &ruleLeafCondition,
			})
		}

		pv := vs["previous"].(*pluginsdk.Set)
		if len(pv.List()) > 0 {
			ruleLeafCondition := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
			for _, e := range pv.List() {
				event := e.(string)
				ruleLeafCondition = append(ruleLeafCondition, activitylogalertsapis.AlertRuleLeafCondition{
					Field:  pointer.To("properties.previousHealthStatus"),
					Equals: pointer.To(event),
				})
			}
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				AnyOf: &ruleLeafCondition,
			})
		}

		rv := vs["reason"].(*pluginsdk.Set)
		if len(rv.List()) > 0 {
			ruleLeafCondition := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
			for _, e := range rv.List() {
				event := e.(string)
				ruleLeafCondition = append(ruleLeafCondition, activitylogalertsapis.AlertRuleLeafCondition{
					Field:  pointer.To("properties.cause"),
					Equals: pointer.To(event),
				})
			}
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				AnyOf: &ruleLeafCondition,
			})
		}
	}
	return conditions
}

func expandServiceHealth(serviceHealth []any, conditions []activitylogalertsapis.AlertRuleAnyOfOrLeafCondition) []activitylogalertsapis.AlertRuleAnyOfOrLeafCondition {
	for _, serviceItem := range serviceHealth {
		if serviceItem == nil {
			continue
		}
		vs := serviceItem.(map[string]any)
		rv := vs["locations"].(*pluginsdk.Set)
		if len(rv.List()) > 0 {
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				Field:       pointer.To("properties.impactedServices[*].ImpactedRegions[*].RegionName"),
				ContainsAny: pluginsdk.ExpandStringSlice(rv.List()),
			})
		}

		ev := vs["events"].(*pluginsdk.Set)
		if len(ev.List()) > 0 {
			ruleLeafCondition := make([]activitylogalertsapis.AlertRuleLeafCondition, 0)
			for _, e := range ev.List() {
				event := e.(string)
				ruleLeafCondition = append(ruleLeafCondition, activitylogalertsapis.AlertRuleLeafCondition{
					Field:  pointer.To("properties.incidentType"),
					Equals: pointer.To(event),
				})
			}
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				AnyOf: &ruleLeafCondition,
			})
		}

		sv := vs["services"].(*pluginsdk.Set)
		if len(sv.List()) > 0 {
			conditions = append(conditions, activitylogalertsapis.AlertRuleAnyOfOrLeafCondition{
				Field:       pointer.To("properties.impactedServices[*].ServiceName"),
				ContainsAny: pluginsdk.ExpandStringSlice(sv.List()),
			})
		}
	}
	return conditions
}

func expandMonitorActivityLogAlertAction(input []any) activitylogalertsapis.ActionList {
	actions := make([]activitylogalertsapis.ActionGroup, 0)
	for _, item := range input {
		v := item.(map[string]any)
		if agID := v["action_group_id"].(string); agID != "" {
			props := make(map[string]string)
			if pVal, ok := v["webhook_properties"]; ok {
				for pk, pv := range pVal.(map[string]any) {
					props[pk] = pv.(string)
				}
			}

			actions = append(actions, activitylogalertsapis.ActionGroup{
				ActionGroupId:     agID,
				WebhookProperties: &props,
			})
		}
	}
	return activitylogalertsapis.ActionList{
		ActionGroups: &actions,
	}
}
