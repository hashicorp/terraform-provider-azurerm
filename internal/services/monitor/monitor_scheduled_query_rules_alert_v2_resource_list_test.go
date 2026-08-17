// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMonitorScheduledQueryRulesAlertV2_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert_v2", "test")
	r := MonitorScheduledQueryRulesAlertV2Resource{}

	listResourceAddress := "azurerm_monitor_scheduled_query_rules_alert_v2.list"
	resourceName := fmt.Sprintf("acctest-sqrv2-1-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctest-rg-monitor-list-%d", data.RandomInteger)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
					querycheck.ExpectIdentity(listResourceAddress, map[string]knownvalue.Check{
						"name":                knownvalue.StringExact(resourceName),
						"resource_group_name": knownvalue.StringExact(resourceGroupName),
						"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (r MonitorScheduledQueryRulesAlertV2Resource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-monitor-list-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctest-ai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_monitor_scheduled_query_rules_alert_v2" "test" {
  count                = 2
  name                 = "acctest-sqrv2-${count.index}-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  evaluation_frequency = "PT5M"
  window_duration      = "PT5M"
  scopes               = [azurerm_application_insights.test.id]
  severity             = 3
  criteria {
    query                   = <<-QUERY
      requests
	    | summarize CountByCountry=count() by client_CountryOrRegion
	  QUERY
    time_aggregation_method = "Count"
    threshold               = 5.0
    operator                = "Equal"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MonitorScheduledQueryRulesAlertV2Resource) basicListQuery() string {
	return `
list "azurerm_monitor_scheduled_query_rules_alert_v2" "list" {
  provider = azurerm
  config {}
}
`
}

func (r MonitorScheduledQueryRulesAlertV2Resource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_monitor_scheduled_query_rules_alert_v2" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctest-rg-monitor-list-%[1]d"
  }
}
`, data.RandomInteger)
}
