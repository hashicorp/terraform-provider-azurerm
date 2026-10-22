// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_activity_log_alert

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMonitorActivityLogAlertFlatten(d *pluginsdk.ResourceData, id *activitylogalertsapis.ActivityLogAlertId, model *activitylogalertsapis.ActivityLogAlertResource) error {
	d.Set("name", id.ActivityLogAlertName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			d.Set("enabled", props.Enabled)
			d.Set("description", props.Description)

			var scopes []any
			if props.Scopes != nil {
				scopes = pluginsdk.FlattenSlice(&props.Scopes)
			}
			if err := d.Set("scopes", scopes); err != nil {
				return fmt.Errorf("setting `scopes`: %+v", err)
			}

			if err := d.Set("criteria", flattenMonitorActivityLogAlertCriteria(props.Condition)); err != nil {
				return fmt.Errorf("setting `criteria`: %+v", err)
			}
			if err := d.Set("action", flattenMonitorActivityLogAlertAction(props.Actions)); err != nil {
				return fmt.Errorf("setting `action`: %+v", err)
			}
		}
		if err := d.Set("tags", pluginsdk.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenMonitorActivityLogAlertCriteria(input activitylogalertsapis.AlertRuleAllOfCondition) []any {
	result := make(map[string]any)
	if input.AllOf == nil {
		return []any{result}
	}
	for _, condition := range input.AllOf {
		if condition.Field != nil && condition.Equals != nil {
			switch strings.ToLower(*condition.Field) {
			case "operationname":
				result["operation_name"] = *condition.Equals
			case "resourceprovider":
				result["resource_provider"] = *condition.Equals
			case "resourcetype":
				result["resource_type"] = *condition.Equals
			case "resourcegroup":
				result["resource_group"] = *condition.Equals
			case "resourceid":
				result["resource_id"] = *condition.Equals
			case "substatus":
				result["sub_status"] = *condition.Equals
			case "properties.recommendationtype":
				result["recommendation_type"] = *condition.Equals
			case "properties.recommendationcategory":
				result["recommendation_category"] = *condition.Equals
			case "properties.recommendationimpact":
				result["recommendation_impact"] = *condition.Equals
			case "caller", "category", "level", "status":
				result[*condition.Field] = *condition.Equals
			}
		}

		if condition.Field != nil && condition.ContainsAny != nil && len(*condition.ContainsAny) > 0 {
			switch strings.ToLower(*condition.Field) {
			case "resourceprovider":
				result["resource_providers"] = *condition.ContainsAny
			case "resourcetype":
				result["resource_types"] = *condition.ContainsAny
			case "resourcegroup":
				result["resource_groups"] = *condition.ContainsAny
			case "resourceid":
				result["resource_ids"] = *condition.ContainsAny
			case "substatus":
				result["sub_statuses"] = *condition.ContainsAny
			case "level":
				result["levels"] = *condition.ContainsAny
			case "status":
				result["statuses"] = *condition.ContainsAny
			}
		}

		if condition.AnyOf != nil && len(*condition.AnyOf) > 0 {
			values := make([]string, 0)
			for _, leafCondition := range *condition.AnyOf {
				if leafCondition.Field != nil && leafCondition.Equals != nil {
					values = append(values, *leafCondition.Equals)
				}
				switch strings.ToLower(*leafCondition.Field) {
				case "resourceprovider":
					result["resource_providers"] = values
				case "resourcetype":
					result["resource_types"] = values
				case "resourcegroup":
					result["resource_groups"] = values
				case "resourceid":
					result["resource_ids"] = values
				case "substatus":
					result["sub_statuses"] = values
				case "level":
					result["levels"] = values
				case "status":
					result["statuses"] = values
				}
			}
		}
	}

	if result["category"] == "ResourceHealth" {
		flattenMonitorActivityLogAlertResourceHealth(input, result)
	}

	if result["category"] == "ServiceHealth" {
		flattenMonitorActivityLogAlertServiceHealth(input, result)
	}

	return []any{result}
}

func flattenMonitorActivityLogAlertResourceHealth(input activitylogalertsapis.AlertRuleAllOfCondition, result map[string]any) {
	rhResult := make(map[string]any)

	for _, condition := range input.AllOf {
		if condition.Field == nil && condition.AnyOf != nil && len(*condition.AnyOf) > 0 {
			values := []string{}
			for _, cond := range *condition.AnyOf {
				if cond.Field != nil && cond.Equals != nil {
					values = append(values, *cond.Equals)
				}
				switch strings.ToLower(*cond.Field) {
				case "properties.currenthealthstatus":
					rhResult["current"] = values
				case "properties.previoushealthstatus":
					rhResult["previous"] = values
				case "properties.cause":
					rhResult["reason"] = values
				}
			}
		}
	}

	result["resource_health"] = []any{rhResult}
}

func flattenMonitorActivityLogAlertServiceHealth(input activitylogalertsapis.AlertRuleAllOfCondition, result map[string]any) {
	shResult := make(map[string]any)
	for _, condition := range input.AllOf {
		if condition.Field != nil && condition.ContainsAny != nil && len(*condition.ContainsAny) > 0 {
			switch strings.ToLower(*condition.Field) {
			case "properties.impactedservices[*].impactedregions[*].regionname":
				shResult["locations"] = *condition.ContainsAny
			case "properties.impactedservices[*].servicename":
				shResult["services"] = *condition.ContainsAny
			}
		}
		if condition.Field == nil && condition.AnyOf != nil && len(*condition.AnyOf) > 0 {
			events := []string{}
			for _, evCond := range *condition.AnyOf {
				if evCond.Field != nil && evCond.Equals != nil {
					events = append(events, *evCond.Equals)
				}
			}
			shResult["events"] = events
		}
	}

	result["service_health"] = []any{shResult}
}

func flattenMonitorActivityLogAlertAction(input activitylogalertsapis.ActionList) (result []any) {
	result = make([]any, 0)
	if input.ActionGroups == nil {
		return
	}
	for _, action := range *input.ActionGroups {
		v := make(map[string]any)

		v["action_group_id"] = action.ActionGroupId

		props := make(map[string]any)
		if action.WebhookProperties != nil {
			for pk, pv := range *action.WebhookProperties {
				props[pk] = pv
			}
		}
		v["webhook_properties"] = props

		result = append(result, v)
	}
	return result
}
