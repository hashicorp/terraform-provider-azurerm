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

type SubscriptionTemplateDeploymentResource struct {
}

func TestAccSubscriptionTemplateDeployment_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.emptyConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// set some tags
			Config: r.emptyWithTagsConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_singleItemUpdatingParams(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

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

func TestAccSubscriptionTemplateDeployment_singleItemUpdatingTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.singleItemWithResourceGroupConfig(data, "first"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.singleItemWithResourceGroupConfig(data, "second"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_withOutputs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withOutputsConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output_content").HasValue("{\"testOutput\":{\"type\":\"String\""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_switchTemplateDeploymentBetweenLinkAndContent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withTemplateLinkAndParametersLinkConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withDeploymentContents(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTemplateLinkAndParametersLinkConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_updateTemplateLinkAndParametersLink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withTemplateLinkAndParametersLinkConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTemplateLinkAndParametersLinkConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_updateExpressionEvaluationOption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withExpressionEvaluationOptionConfig(data, "Inner"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// The Azure API doesn't return the `expression_evaluation_option` property in the response.
		// Bug: https://github.com/Azure/azure-rest-api-specs/issues/12326
		data.ImportStep("expression_evaluation_option"),
		{
			Config: r.withExpressionEvaluationOptionConfig(data, "Outer"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("expression_evaluation_option"),
	})
}

func (t SubscriptionTemplateDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SubscriptionTemplateDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Resource.DeploymentsClient.GetAtSubscriptionScope(ctx, id.DeploymentName)
	if err != nil {
		return nil, fmt.Errorf("reading Subscription Template Deployment (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (SubscriptionTemplateDeploymentResource) emptyConfig(data acceptance.TestData) string {
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

func (SubscriptionTemplateDeploymentResource) emptyWithTagsConfig(data acceptance.TestData) string {
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

func (SubscriptionTemplateDeploymentResource) singleItemWithParameterConfig(data acceptance.TestData, value string) string {
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

func (SubscriptionTemplateDeploymentResource) singleItemWithResourceGroupConfig(data acceptance.TestData, tagValue string) string {
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

func (SubscriptionTemplateDeploymentResource) withOutputsConfig(data acceptance.TestData) string {
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

func (SubscriptionTemplateDeploymentResource) withTemplateLinkAndParametersLinkConfig(data acceptance.TestData) string {
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

func (SubscriptionTemplateDeploymentResource) updateTemplateLinkAndParametersLinkConfig(data acceptance.TestData) string {
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

func (SubscriptionTemplateDeploymentResource) withExpressionEvaluationOptionConfig(data acceptance.TestData, scope string) string {
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

func (SubscriptionTemplateDeploymentResource) withDeploymentContents(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_subscription_template_deployment" "test" {
  name               = "acctest-SubDeploy-%d"
  location           = "%s"
  template_content   = <<TEMPLATE
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
