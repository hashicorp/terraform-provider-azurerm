// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_elasticsearch

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	elasticClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceElasticsearchUpdate(d *pluginsdk.ResourceData, meta any) error {
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := monitorsresource.ParseMonitorID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("logs") {
		client := elasticClient.From(meta).TagRuleClient
		tagRuleId := rules.NewTagRuleID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, "default")
		tagRule := expandTagRule(d.Get("logs").([]any))
		body := rules.MonitoringTagRules{
			Properties: &rules.MonitoringTagRulesProperties{
				LogRules: tagRule,
			},
		}
		if _, err := client.TagRulesCreateOrUpdate(ctx, tagRuleId, body); err != nil {
			return fmt.Errorf("updating `logs` from %s: %+v", *id, err)
		}
	}

	if d.HasChange("tags") {
		client := elasticClient.From(meta).MonitorClient
		body := monitorsresource.ElasticMonitorResourceUpdateParameters{
			Tags: tags.Expand(d.Get("tags").(map[string]any)),
		}
		if _, err := client.MonitorsUpdate(ctx, *id, body); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	return resourceElasticsearchRead(d, meta)
}
