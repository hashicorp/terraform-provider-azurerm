// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_prometheus_rule_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroupresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AlertPrometheusRuleGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := prometheusrulegroupresources.NewPrometheusRuleGroupID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.PrometheusRuleGroupsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := prometheusrulegroupresources.PrometheusRuleGroupResource{
				Location: location.Normalize(model.Location),
				Properties: prometheusrulegroupresources.PrometheusRuleGroupProperties{
					Enabled: pointer.To(model.RuleGroupEnabled),
					Scopes:  model.Scopes,
				},
				Tags: pointer.To(model.Tags),
			}

			properties.Properties.ClusterName = pointer.To(model.ClusterName)
			properties.Properties.Description = pointer.To(model.Description)
			if _, ok := metadata.ResourceData.GetOk("interval"); ok {
				properties.Properties.Interval = pointer.To(model.Interval)
			}
			properties.Properties.Rules = expandPrometheusRuleModel(model.Rule, metadata.ResourceData)

			if _, err := client.PrometheusRuleGroupsCreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandPrometheusRuleModel(inputList []PrometheusRuleModel, d *schema.ResourceData) []prometheusrulegroupresources.PrometheusRule {
	outputList := make([]prometheusrulegroupresources.PrometheusRule, 0)

	for i, v := range inputList {
		output := prometheusrulegroupresources.PrometheusRule{
			Enabled:    pointer.To(v.Enabled),
			Expression: v.Expression,
			Labels:     pointer.To(v.Labels),
		}

		if v.Alert != "" {
			output.Actions = expandPrometheusRuleGroupActionModel(v.Action)
			output.Alert = pointer.To(v.Alert)
			if v, ok := d.GetOk(fmt.Sprintf("rule.%d.severity", i)); ok {
				output.Severity = pointer.To(int64(v.(int)))
			}
			output.Annotations = pointer.To(v.Annotations)
			output.ResolveConfiguration = expandPrometheusRuleAlertResolutionModel(v.AlertResolution)

			if v.For != "" {
				output.For = pointer.To(v.For)
			}
		} else {
			// action, alert, severity, annotations, for, alert_resolution must be empty when type is recording rule
			output.Record = pointer.To(v.Record)
		}
		outputList = append(outputList, output)
	}

	return outputList
}

func expandPrometheusRuleGroupActionModel(inputList []PrometheusRuleGroupActionModel) *[]prometheusrulegroupresources.PrometheusRuleGroupAction {
	outputList := make([]prometheusrulegroupresources.PrometheusRuleGroupAction, 0)
	for _, v := range inputList {
		output := prometheusrulegroupresources.PrometheusRuleGroupAction{
			ActionProperties: pointer.To(v.ActionProperties),
		}
		output.ActionGroupId = pointer.To(v.ActionGroupId)
		outputList = append(outputList, output)
	}

	return &outputList
}

func expandPrometheusRuleAlertResolutionModel(inputList []PrometheusRuleAlertResolutionModel) *prometheusrulegroupresources.PrometheusRuleResolveConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := prometheusrulegroupresources.PrometheusRuleResolveConfiguration{
		AutoResolved: pointer.To(input.AutoResolved),
	}
	output.TimeToResolve = pointer.To(input.TimeToResolve)

	return &output
}
