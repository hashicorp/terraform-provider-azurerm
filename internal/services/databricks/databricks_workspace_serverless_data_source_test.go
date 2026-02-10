// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DatabricksWorkspaceServerlessDataSource struct{}

func TestAccDatabricksWorkspaceServerlessDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("workspace_id").Exists(),
				acceptance.TestMatchResourceAttr(data.ResourceName, "workspace_url", regexp.MustCompile("azuredatabricks.net")),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("Sandbox"),
				check.That(data.ResourceName).Key("tags.Label").HasValue("Test"),
			),
		},
	})
}

func TestAccDatabricksWorkspaceServerlessDataSource_enhancedComplianceSecurity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databricks_workspace_serverless", "test")
	r := DatabricksWorkspaceServerlessDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.enhancedSecurityCompliance(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("enhanced_security_compliance.#").HasValue("1"),
				check.That(data.ResourceName).Key("enhanced_security_compliance.0.automatic_cluster_update_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("enhanced_security_compliance.0.compliance_security_profile_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("enhanced_security_compliance.0.compliance_security_profile_standards.#").HasValue("1"),
				check.That(data.ResourceName).Key("enhanced_security_compliance.0.enhanced_security_monitoring_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("workspace_id").Exists(),
				acceptance.TestMatchResourceAttr(data.ResourceName, "workspace_url", regexp.MustCompile("azuredatabricks.net")),
			),
		},
	})
}

func (DatabricksWorkspaceServerlessDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-databricks-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DatabricksWorkspaceServerlessDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_workspace_serverless" "test" {
  name                = "acctest-dbsw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    Environment = "Sandbox"
    Label       = "Test"
  }
}

data "azurerm_databricks_workspace_serverless" "test" {
  name                = azurerm_databricks_workspace_serverless.test.name
  resource_group_name = azurerm_databricks_workspace_serverless.test.resource_group_name
}
`, r.template(data), data.RandomInteger)
}

func (r DatabricksWorkspaceServerlessDataSource) enhancedSecurityCompliance(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_workspace_serverless" "test" {
  name                = "acctest-dbsw-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  enhanced_security_compliance {
    automatic_cluster_update_enabled      = true
    compliance_security_profile_enabled   = true
    compliance_security_profile_standards = ["HIPAA"]
    enhanced_security_monitoring_enabled  = true
  }
}

data "azurerm_databricks_workspace_serverless" "test" {
  name                = azurerm_databricks_workspace_serverless.test.name
  resource_group_name = azurerm_databricks_workspace_serverless.test.resource_group_name
}
`, r.template(data), data.RandomInteger)
}
