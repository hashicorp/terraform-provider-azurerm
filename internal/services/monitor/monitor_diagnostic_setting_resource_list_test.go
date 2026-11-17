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

func TestAccMonitorDiagnosticSetting_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "list")
	r := MonitorDiagnosticSettingResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_monitor_diagnostic_setting.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_monitor_diagnostic_setting.list",
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

func (MonitorDiagnosticSettingResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_management_group" "test" {
  name = "acctestMG-%[1]d"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_monitor_diagnostic_setting" "test1" {
  name                       = "acctest-MDS1-%[1]d"
  target_resource_id         = azurerm_management_group.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category = "Administrative"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test2" {
  name                       = "acctest-MDS2-%[1]d"
  target_resource_id         = azurerm_management_group.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category = "Policy"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MonitorDiagnosticSettingResource) basicListQuery() string {
	return `
list "azurerm_monitor_diagnostic_setting" "list" {
  provider = azurerm
  config {
    target_resource_id = azurerm_management_group.test.id
  }
}
`
}
