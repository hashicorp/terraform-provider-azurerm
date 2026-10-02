// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccResourceGroupTemplateDeployment_V0ToV1_501 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.0.1 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccResourceGroupTemplateDeployment_V0ToV1_501(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyConfig(data, "Complete"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Resources/deployments/acctest", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.emptyV0(data, "Complete"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Ensure import is a no-op to prevent the resource from normalizing the ID due to a combined CreateUpdate func
					plancheck.ExpectResourceAction(importedResourceName, plancheck.ResourceActionNoop),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.resources/deployments/acctest", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.emptyImported(data, "Complete"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Resources/deployments/acctest", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.0.1")
}

func (ResourceGroupTemplateDeploymentResource) emptyV0(data acceptance.TestData, deploymentMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = %[2]q
}

resource "azurerm_resource_group_template_deployment" "test-import" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  deployment_mode     = %[3]q

  template_content = <<TEMPLATE
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {},
  "variables": {},
  "resources": []
}
TEMPLATE

  lifecycle {
    # ignore template_content, the template exported by the API on import may be normalised differently to the config
    ignore_changes = [template_content]
  }
}

removed {
  from = azurerm_resource_group_template_deployment.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_resource_group_template_deployment.test-import
  id = "/subscriptions/%[4]s/resourcegroups/acctestRG-%[1]d/providers/microsoft.resources/deployments/acctest"
}
`, data.RandomInteger, data.Locations.Primary, deploymentMode, data.Subscriptions.Primary)
}

func (ResourceGroupTemplateDeploymentResource) emptyImported(data acceptance.TestData, deploymentMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = %[2]q
}

resource "azurerm_resource_group_template_deployment" "test-import" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  deployment_mode     = %[3]q

  template_content = <<TEMPLATE
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {},
  "variables": {},
  "resources": []
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary, deploymentMode)
}
