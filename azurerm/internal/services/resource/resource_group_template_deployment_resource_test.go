package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ResourceGroupTemplateDeploymentResource struct {
}

func TestAccResourceGroupTemplateDeployment_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.emptyConfig(data, "Complete"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// set some tags
			Config: r.emptyWithTagsConfig(data, "Complete"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_incremental(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.emptyConfig(data, "Incremental"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// set some tags
			Config: r.emptyWithTagsConfig(data, "Incremental"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_singleItemIncorrectCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.singleItemWithIncorrectCasingConfig(data, "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_singleItemUpdateTemplateWithUnchangedParams(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.singleItemWithParameterConfigAndVariable(data, "templateVariableValue1", "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Update the ARM template to trigger update on resource BUT don't vary the parameter value
			// Validates fix for bug #8840
			Config: r.singleItemWithParameterConfigAndVariable(data, "forceupdatetemplate", "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_singleItemUpdatingParams(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.singleItemWithParameterConfig(data, "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.singleItemWithParameterConfig(data, "second"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_singleItemUpdatingTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.singleItemWithPublicIPConfig(data, "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.singleItemWithPublicIPConfig(data, "second"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_withOutputs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withOutputsConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output_content").HasValue("{\"testOutput\":{\"type\":\"String\",\"value\":\"some-value\"}}"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_multipleItems(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleItemsConfig(data, "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleItemsConfig(data, "second"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_multipleNestedItems(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleNestedItemsConfig(data, "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleNestedItemsConfig(data, "second"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_childItems(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.childItemsConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.childItemsConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_templateSpecEmpty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.templateSpecVersionConfigEmpty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupTemplateDeployment_templateSpecResources(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_template_deployment", "test")
	r := ResourceGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.templateSpecVersionConfigResources(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t ResourceGroupTemplateDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ResourceGroupTemplateDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Resource.DeploymentsClient.Get(ctx, id.ResourceGroup, id.DeploymentName)
	if err != nil {
		return nil, fmt.Errorf("reading Management Lock (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ResourceGroupTemplateDeploymentResource) emptyConfig(data acceptance.TestData, deploymentMode string) string {
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
  deployment_mode     = %q

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

func (ResourceGroupTemplateDeploymentResource) templateSpecVersionConfigEmpty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = %q
}

data "azurerm_template_spec_version" "test" {
  name                = "acctest-standing-data-empty"
  resource_group_name = "standing-data-for-acctest"
  version             = "v1.0.0"
}

resource "azurerm_resource_group_template_deployment" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  deployment_mode     = "Incremental"

  template_spec_version_id = data.azurerm_template_spec_version.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ResourceGroupTemplateDeploymentResource) templateSpecVersionConfigResources(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = %q
}

data "azurerm_template_spec_version" "test" {
  name                = "acctest-standing-data-for-rg"
  resource_group_name = "standing-data-for-acctest"
  version             = "v1.0.0"
}

resource "azurerm_resource_group_template_deployment" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  deployment_mode     = "Incremental"

  template_spec_version_id = data.azurerm_template_spec_version.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ResourceGroupTemplateDeploymentResource) emptyWithTagsConfig(data acceptance.TestData, deploymentMode string) string {
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
  deployment_mode     = %q
  tags = {
    Hello = "World"
  }

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

func (ResourceGroupTemplateDeploymentResource) singleItemWithIncorrectCasingConfig(data acceptance.TestData, value string) string {
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
  "parameters": {
    "someParam": {
      "type": "String",
      "allowedValues": [
        "first",
        "second",
        "third"
      ]
    }
  },
  "variables": {},
  "resources": [
    {
      "type": "microsoft.insights/actionGroups",
      "apiVersion": "2019-06-01",
      "name": "acctestTemplateDeployAG-%d",
      "location": "Global",
      "dependsOn": [],
      "tags": {},
      "properties": {
        "groupShortName": "rick-c137",
        "enabled": true,
        "emailReceivers": [
          {
            "name": "Rick Sanchez",
            "emailAddress": "rick@example.com"
          }
        ],
        "smsReceivers": [],
        "webhookReceivers": []
      }
    }
  ]
}
TEMPLATE

  parameters_content = <<PARAM
{
  "someParam": {
   "value": %q
  }
}
PARAM
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, value)
}

func (ResourceGroupTemplateDeploymentResource) singleItemWithParameterConfigAndVariable(data acceptance.TestData, templateVar string, paramVal string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
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
  "parameters": {
    "someParam": {
      "type": "String",
      "allowedValues": [
        "first",
        "second",
        "third"
      ]
    }
  },
  "variables": {
    "forceChangeInTemplate": %q
  },
  "resources": []
}
TEMPLATE

  parameters_content = <<PARAM
{
  "someParam": {
   "value": %q
  }
}
PARAM
}
`, data.RandomInteger, data.Locations.Primary, templateVar, paramVal)
}

func (ResourceGroupTemplateDeploymentResource) singleItemWithParameterConfig(data acceptance.TestData, value string) string {
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
  "parameters": {
    "someParam": {
      "type": "String",
      "allowedValues": [
        "first",
        "second",
        "third"
      ]
    }
  },
  "variables": {},
  "resources": []
}
TEMPLATE

  parameters_content = <<PARAM
{
  "someParam": {
   "value": %q
  }
}
PARAM
}
`, data.RandomInteger, data.Locations.Primary, value)
}

func (ResourceGroupTemplateDeploymentResource) singleItemWithPublicIPConfig(data acceptance.TestData, tagValue string) string {
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
  "resources": [
    {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "2015-06-15",
      "name": "acctestpip-%d",
      "location": "[resourceGroup().location]",
      "properties": {
        "publicIPAllocationMethod": "Dynamic"
      },
      "tags": {
        "Hello": %q
      }
    }
  ]
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, tagValue)
}

func (ResourceGroupTemplateDeploymentResource) withOutputsConfig(data acceptance.TestData) string {
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

func (ResourceGroupTemplateDeploymentResource) multipleItemsConfig(data acceptance.TestData, value string) string {
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
  "resources": [
    {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "2015-06-15",
      "name": "acctestpip-1-%d",
      "location": "[resourceGroup().location]",
      "properties": {
        "publicIPAllocationMethod": "Dynamic"
      },
      "tags": {
        "Hello": %q
      }
    },
    {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "2015-06-15",
      "name": "acctestpip-2-%d",
      "location": "[resourceGroup().location]",
      "properties": {
        "publicIPAllocationMethod": "Dynamic"
      },
      "tags": {
        "Hello": %q
      }
    }
  ]
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, value, data.RandomInteger, value)
}

func (ResourceGroupTemplateDeploymentResource) multipleNestedItemsConfig(data acceptance.TestData, value string) string {
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
  "resources": [
    {
      "type": "Microsoft.Network/virtualNetworks",
      "apiVersion": "2020-05-01",
      "name": "parent-network",
      "location": "[resourceGroup().location]",
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            "10.0.0.0/16"
          ]
        }
      },
      "tags": {
        "Hello": %q
      },
      "resources": [
        {
          "type": "subnets",
          "apiVersion": "2020-05-01",
          "location": "[resourceGroup().location]",
          "name": "first",
          "dependsOn": [
            "parent-network"
          ],
          "properties": {
            "addressPrefix": "10.0.1.0/24"
          }
        },
        {
          "type": "subnets",
          "apiVersion": "2020-05-01",
          "location": "[resourceGroup().location]",
          "name": "second",
          "dependsOn": [
            "parent-network",
            "first"
          ],
          "properties": {
            "addressPrefix": "10.0.2.0/24"
          }
        }
      ]
    }
  ]
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary, value)
}

func (ResourceGroupTemplateDeploymentResource) childItemsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = %q
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_resource_group_template_deployment" "test" {
  name                = "acctest"
  resource_group_name = azurerm_resource_group.test.name
  deployment_mode     = "Incremental"

  template_content = <<TEMPLATE
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {},
  "variables": {},
  "resources": [
    {
      "type": "Microsoft.Network/routeTables/routes",
      "apiVersion": "2020-06-01",
      "name": "${azurerm_route_table.test.name}/child-route",
      "location": "[resourceGroup().location]",
      "properties": {
        "addressPrefix": "10.2.0.0/16",
        "nextHopType": "none"
      }
    }
  ]
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
