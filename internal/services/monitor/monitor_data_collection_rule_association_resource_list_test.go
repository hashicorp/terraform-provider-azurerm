// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMonitorDataCollectionRuleAssociation_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_data_collection_rule_association", "list")
	r := MonitorDataCollectionRuleAssociationResource{}

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
				Config: r.basicListQueryTargetResource(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_monitor_data_collection_rule_association.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_monitor_data_collection_rule_association.list",
						map[string]knownvalue.Check{
							"name":         knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_uri": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryRule(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_monitor_data_collection_rule_association.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_monitor_data_collection_rule_association.list",
						map[string]knownvalue.Check{
							"name":         knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_uri": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
		},
	})
}

func (r MonitorDataCollectionRuleAssociationResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_monitor_data_collection_rule" "test2" {
  name                = "acctestmdcr2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  destinations {
    azure_monitor_metrics {
      name = "test-destination-metrics"
    }
  }
  data_flow {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["test-destination-metrics"]
  }
}

resource "azurerm_monitor_data_collection_rule_association" "test1" {
  name                    = "test-dcra-1-%[2]d"
  target_resource_id      = azurerm_linux_virtual_machine.test.id
  data_collection_rule_id = azurerm_monitor_data_collection_rule.test.id
}

resource "azurerm_monitor_data_collection_rule_association" "test2" {
  name                    = "test-dcra-2-%[2]d"
  target_resource_id      = azurerm_linux_virtual_machine.test.id
  data_collection_rule_id = azurerm_monitor_data_collection_rule.test2.id
}
`, r.template(data), data.RandomInteger)
}

func (MonitorDataCollectionRuleAssociationResource) basicListQueryTargetResource() string {
	return `
list "azurerm_monitor_data_collection_rule_association" "list" {
  provider = azurerm
  config {
    target_resource_id = azurerm_linux_virtual_machine.test.id
  }
}
`
}

func (MonitorDataCollectionRuleAssociationResource) basicListQueryRule() string {
	return `
list "azurerm_monitor_data_collection_rule_association" "list" {
  provider = azurerm
  config {
    data_collection_rule_id = azurerm_monitor_data_collection_rule.test.id
  }
}
`
}
