// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_elasticsearch

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceElasticsearchCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Elastic.MonitorClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := monitorsresource.NewMonitorID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.MonitorsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_elastic_cloud_elasticsearch", id.ID())
		}
	}

	monitoringStatus := monitorsresource.MonitoringStatusDisabled
	if d.Get("monitoring_enabled").(bool) {
		monitoringStatus = monitorsresource.MonitoringStatusEnabled
	}

	body := monitorsresource.ElasticMonitorResource{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &monitorsresource.MonitorProperties{
			MonitoringStatus: &monitoringStatus,
			UserInfo: &monitorsresource.UserInfo{
				EmailAddress: pointer.To(d.Get("elastic_cloud_email_address").(string)),
			},
		},
		Sku: &monitorsresource.ResourceSku{
			Name: d.Get("sku_name").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.MonitorsCreateCallbackThenPoll(ctx, id, body, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if v, ok := d.GetOk("logs"); ok {
		tagRulesClient := meta.(*clients.Client).Elastic.TagRuleClient
		tagRuleId := rules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
		tagRule := rules.MonitoringTagRules{
			Properties: &rules.MonitoringTagRulesProperties{
				LogRules: expandTagRule(v.([]any)),
			},
		}
		if _, err := tagRulesClient.TagRulesCreateOrUpdate(ctx, tagRuleId, tagRule); err != nil {
			return fmt.Errorf("updating the logs for %s: %+v", id, err)
		}
	}

	return resourceElasticsearchRead(d, meta)
}

func expandTagRule(input []any) *rules.LogRules {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]any)
	filteringTags := make([]rules.FilteringTag, 0)
	for _, v := range raw["filtering_tag"].([]any) {
		item := v.(map[string]any)

		filteringTags = append(filteringTags, rules.FilteringTag{
			Action: pointer.ToEnum[rules.TagAction](item["action"].(string)),
			Name:   pointer.To(item["name"].(string)),
			Value:  pointer.To(item["value"].(string)),
		})
	}

	sendAzureAdLogs := raw["send_azuread_logs"].(bool)
	sendActivityLogs := raw["send_activity_logs"].(bool)
	sendSubscriptionLogs := raw["send_subscription_logs"].(bool)

	return &rules.LogRules{
		FilteringTags:        &filteringTags,
		SendAadLogs:          pointer.To(sendAzureAdLogs),
		SendActivityLogs:     pointer.To(sendActivityLogs),
		SendSubscriptionLogs: pointer.To(sendSubscriptionLogs),
	}
}
