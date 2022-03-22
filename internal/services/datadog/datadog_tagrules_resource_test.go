package datadog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TagRulesDatadogMonitorResource struct{}

func TestAccDatadogMonitorTagRules_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_tagrules", "test")
	r := TagRulesDatadogMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatadogMonitorTagRules_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_tagrules", "test")
	r := TagRulesDatadogMonitorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r TagRulesDatadogMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DatadogTagRulesID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Datadog.TagRulesClient.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Datadog Monitor %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	if (*resp.Properties.LogRules.SendResourceLogs == false && *resp.Properties.LogRules.SendSubscriptionLogs == false) || *resp.Properties.MetricRules.FilteringTags == nil {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r TagRulesDatadogMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
	provider "azurerm" {
	features {}
	}

	data "azurerm_resource_group" "test" {
	name     = "acctest-datadog"
	}
	`)
}

func (r TagRulesDatadogMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	
	resource "azurerm_datadog_monitor_tagrules" "test" {
		name = "test-terraform-acctests"
		resource_group_name = "acctest-datadog"
		log_rules{
			send_subscription_logs = true
		}
		metric_rules{
			filtering_tag {
				name = "Vidhi"
				value = "Testing-Logs"
				action = "Include"
			}
		}
	}
`)
}

func (r TagRulesDatadogMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`

	resource "azurerm_datadog_monitor_tagrules" "test" {
		name = "test-terraform-acctests"
		resource_group_name = "acctest-datadog"
		log_rules{
			send_subscription_logs = false
			send_resource_logs = true
			filtering_tag {
				name = "Vidhi"
				value = "Testing-Logs"
				action = "Include"
			}
		}
		metric_rules{
			filtering_tag {
				name = "Vidhi"
				value = "Testing-Logs"
				action = "Include"
			}
		}
	}
`)
}
