package resource_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccSubscriptionTemplateDeployment_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSubscriptionTemplateDeploymentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: subscriptionTemplateDeployment_emptyConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// set some tags
				Config: subscriptionTemplateDeployment_emptyWithTagsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccSubscriptionTemplateDeployment_singleItemUpdatingParams(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSubscriptionTemplateDeploymentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: subscriptionTemplateDeployment_singleItemWithParameterConfig(data, "first"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: subscriptionTemplateDeployment_singleItemWithParameterConfig(data, "second"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccSubscriptionTemplateDeployment_singleItemUpdatingTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSubscriptionTemplateDeploymentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: subscriptionTemplateDeployment_singleItemWithResourceGroupConfig(data, "first"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: subscriptionTemplateDeployment_singleItemWithResourceGroupConfig(data, "second"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccSubscriptionTemplateDeployment_withOutputs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSubscriptionTemplateDeploymentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: subscriptionTemplateDeployment_withOutputsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "output_content", "{\"testOutput\":{\"type\":\"String\",\"value\":\"some-value\"}}"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccSubscriptionTemplateDeployment_switchTemplateDeploymentBetweenLinkAndContent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSubscriptionTemplateDeploymentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: subscriptionTemplateDeployment_withTemplateLinkAndParametersLinkConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: subscriptionTemplateDeployment_withDeploymentContents(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: subscriptionTemplateDeployment_withTemplateLinkAndParametersLinkConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccSubscriptionTemplateDeployment_updateTemplateLinkAndParametersLink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSubscriptionTemplateDeploymentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: subscriptionTemplateDeployment_withTemplateLinkAndParametersLinkConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: subscriptionTemplateDeployment_updateTemplateLinkAndParametersLinkConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccSubscriptionTemplateDeployment_updateExpressionEvaluationOption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSubscriptionTemplateDeploymentDestroyed,
		Steps: []resource.TestStep{
			{
				Config: subscriptionTemplateDeployment_withExpressionEvaluationOptionConfig(data, "Inner"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			// The Azure API doesn't return the `expression_evaluation_option` property in the response.
			// Bug: https://github.com/Azure/azure-rest-api-specs/issues/12326
			data.ImportStep("expression_evaluation_option"),
			{
				Config: subscriptionTemplateDeployment_withExpressionEvaluationOptionConfig(data, "Outer"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSubscriptionTemplateDeploymentExists(data.ResourceName),
				),
			},
			data.ImportStep("expression_evaluation_option"),
		},
	})
}

func testCheckSubscriptionTemplateDeploymentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.DeploymentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resp, err := client.GetAtSubscriptionScope(ctx, name)
		if err != nil {
			return fmt.Errorf("bad: Get on deploymentsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("bad: Subscription Template Deployment %q does not exist", name)
		}

		return nil
	}
}

func testCheckSubscriptionTemplateDeploymentDestroyed(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.DeploymentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_subscription_template_deployment" {
			continue
		}

		name := rs.Primary.Attributes["name"]

		resp, err := client.GetAtSubscriptionScope(ctx, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Subscription Template Deployment still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func subscriptionTemplateDeployment_emptyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctestsubdeploy-%d"
  location = %q

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
`, data.RandomInteger, data.Locations.Primary)
}

func subscriptionTemplateDeployment_emptyWithTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctestsubdeploy-%d"
  location = %q
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
`, data.RandomInteger, data.Locations.Primary)
}

func subscriptionTemplateDeployment_singleItemWithParameterConfig(data acceptance.TestData, value string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctestsubdeploy-%d"
  location = %q

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

func subscriptionTemplateDeployment_singleItemWithResourceGroupConfig(data acceptance.TestData, tagValue string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctestsubdeploy-%d"
  location = %q

  template_content = <<TEMPLATE
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {},
  "variables": {},
  "resources": [
    {
      "type": "Microsoft.Resources/resourceGroups",
      "apiVersion": "2018-05-01",
      "location": "%s",
      "name": "acctestrg-%d",
      "properties": {},
      "tags": {
        "Hello": %q
      }
    }
  ]
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Primary, data.RandomInteger, tagValue)
}

func subscriptionTemplateDeployment_withOutputsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctestsubdeploy-%d"
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

func subscriptionTemplateDeployment_withTemplateLinkAndParametersLinkConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctest-SubDeploy-%d"
  location = "%s"

  template_link {
    uri             = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"
    content_version = "1.0.0.0"
  }

  parameters_link {
    uri             = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.parameters.json"
    content_version = "1.0.0.0"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func subscriptionTemplateDeployment_updateTemplateLinkAndParametersLinkConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctest-SubDeploy-%d"
  location = "%s"

  template_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/subscription-deployments/create-rg/azuredeploy.json"
  }

  parameters_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/subscription-deployments/create-rg/azuredeploy.parameters.json"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func subscriptionTemplateDeployment_withExpressionEvaluationOptionConfig(data acceptance.TestData, scope string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctest-SubDeploy-%d"
  location = "%s"

  template_link {
    uri             = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.json"
    content_version = "1.0.0.0"
  }

  parameters_link {
    uri             = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/100-blank-template/azuredeploy.parameters.json"
    content_version = "1.0.0.0"
  }

  expression_evaluation_option {
    scope = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, scope)
}

func subscriptionTemplateDeployment_withDeploymentContents(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctest-SubDeploy-%d"
  location = "%s"

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
   "value": "first"
  }
}
PARAM
}
`, data.RandomInteger, data.Locations.Primary)
}
