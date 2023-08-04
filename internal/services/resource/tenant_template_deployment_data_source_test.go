// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type TenantTemplateDeploymentDataSource struct{}

func TestAccDataSourceTenantTemplateDeployment(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_tenant_template_deployment", "test")
	if data.Client().IsServicePrincipal {
		t.Skip("Skipping due to permissions unavailable on tenant scope")
	}

	r := TenantTemplateDeploymentDataSource{}

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

func (TenantTemplateDeploymentDataSource) withOutputsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_tenant_template_deployment" "test" {
  name     = "acctestTenantDeploy-%d"
  location = %q

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

func (r TenantTemplateDeploymentDataSource) withDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_tenant_template_deployment" "test" {
  name = azurerm_tenant_template_deployment.test.name
}
`, r.withOutputsConfig(data))
}
