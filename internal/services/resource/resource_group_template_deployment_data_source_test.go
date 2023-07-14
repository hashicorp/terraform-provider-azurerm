// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ResourceGroupTemplateDeploymentDataSource struct{}

func TestAccDataSourceResourceGroupTemplateDeployment(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withOutputsConfig(data),
		},
		{
			Config: r.withDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output_content").HasValue("{\"testOutput\":{\"type\":\"String\",\"value\":\"some-value\"}}"),
			),
		},
	})
}

func (ResourceGroupTemplateDeploymentDataSource) withOutputsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = %q
}

resource "azurerm_resource_group_template_deployment" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  deployment_mode     = "Complete"

  template_content = <<TEMPLATE
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {},
  "variables": {},
  "resources": [],
  "outputs": {
    "testOutput": {
      "type": "String",
      "value": "some-value"
    }
  }
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceGroupTemplateDeploymentDataSource) withDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_resource_group_template_deployment" "test" {
  name                = azurerm_resource_group_template_deployment.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, r.withOutputsConfig(data))
}
