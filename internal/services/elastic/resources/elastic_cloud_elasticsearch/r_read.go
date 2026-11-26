// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_elasticsearch

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceElasticsearchRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Elastic.MonitorClient
	logsClient := meta.(*clients.Client).Elastic.TagRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := monitorsresource.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	tagRuleId := rules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
	rulesResp, err := logsClient.TagRulesGet(ctx, tagRuleId)
	if err != nil {
		if !response.WasNotFound(rulesResp.HttpResponse) {
			return fmt.Errorf("retrieving logs for %s: %+v", *id, err)
		}
	}

	d.Set("name", id.MonitorName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			monitoringEnabled := false
			if props.MonitoringStatus != nil {
				monitoringEnabled = *props.MonitoringStatus == monitorsresource.MonitoringStatusEnabled
			}
			d.Set("monitoring_enabled", monitoringEnabled)

			if elastic := props.ElasticProperties; elastic != nil {
				if elastic.ElasticCloudDeployment != nil {
					// AzureSubscriptionId is the same as the subscription deployed into, so no point exposing it
					// ElasticsearchRegion is `{Cloud}-{Region}` - so the same as location/not worth exposing for now?
					d.Set("elastic_cloud_deployment_id", elastic.ElasticCloudDeployment.DeploymentId)
					d.Set("elasticsearch_service_url", elastic.ElasticCloudDeployment.ElasticsearchServiceURL)
					d.Set("kibana_service_url", elastic.ElasticCloudDeployment.KibanaServiceURL)
					d.Set("kibana_sso_uri", elastic.ElasticCloudDeployment.KibanaSsoURL)
				}
				if elastic.ElasticCloudUser != nil {
					d.Set("elastic_cloud_user_id", elastic.ElasticCloudUser.Id)
					d.Set("elastic_cloud_email_address", elastic.ElasticCloudUser.EmailAddress)
					d.Set("elastic_cloud_sso_default_url", elastic.ElasticCloudUser.ElasticCloudSsoDefaultURL)
				}
			}
		}

		skuName := ""
		if model.Sku != nil {
			skuName = model.Sku.Name
		}
		d.Set("sku_name", skuName)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	if err := d.Set("logs", flattenTagRule(rulesResp.Model)); err != nil {
		return fmt.Errorf("setting `logs`: %+v", err)
	}

	return nil
}

func flattenTagRule(input *rules.MonitoringTagRules) []any {
	if input == nil || input.Properties == nil || input.Properties.LogRules == nil {
		return []any{}
	}

	rules := input.Properties.LogRules

	filteringTags := make([]any, 0)
	if rules.FilteringTags != nil {
		for _, v := range *rules.FilteringTags {
			action := ""
			if v.Action != nil {
				action = string(*v.Action)
			}
			filteringTags = append(filteringTags, map[string]any{
				"action": action,
				"name":   pointer.From(v.Name),
				"value":  pointer.From(v.Value),
			})
		}
	}

	return []any{
		map[string]any{
			"filtering_tag":          filteringTags,
			"send_activity_logs":     pointer.From(rules.SendActivityLogs),
			"send_azuread_logs":      pointer.From(rules.SendAadLogs),
			"send_subscription_logs": pointer.From(rules.SendSubscriptionLogs),
		},
	}
}
