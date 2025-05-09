// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubscriptionTemplateDeploymentResource struct{}

func TestAccSubscriptionTemplateDeployment_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// set some tags
			Config: r.emptyWithTagsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_templateSpecVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.templateSpecVersionConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_singleItemUpdatingParams(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleItemWithParameterConfig(data, "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.singleItemWithParameterConfig(data, "second"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_singleItemUpdateTemplateWithUnchangedParams(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleItemWithParameterConfigAndVariable(data, "templateVariableValue1", "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.singleItemWithParameterConfigAndVariable(data, "forceupdatetemplate", "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_singleItemUpdatingTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleItemWithResourceGroupConfig(data, "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.singleItemWithResourceGroupConfig(data, "second"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionTemplateDeployment_withOutputs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_template_deployment", "test")
	r := SubscriptionTemplateDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withOutputsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("output_content").HasValue("{\"testOutput\":{\"type\":\"String\",\"value\":\"some-value\"}}"),
			),
		},
		data.ImportStep(),
	})
}

func (t SubscriptionTemplateDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SubscriptionTemplateDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Resource.LegacyDeploymentsClient.GetAtSubscriptionScope(ctx, id.DeploymentName)
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

func (SubscriptionTemplateDeploymentResource) templateSpecVersionConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_template_spec_version" "test" {
  name                = "acctest-standing-data-for-sub"
  resource_group_name = "standing-data-for-acctest"
  version             = "v1.0.0"
}

resource "azurerm_subscription_template_deployment" "test" {
  name     = "acctestsubdeploy-%d"
  location = %q

  template_spec_version_id = data.azurerm_template_spec_version.test.id

  parameters_content = <<PARAM
{
  "rgName": {
   "value": "acctest-rg-tspec-%d"
  },
  "rgLocation": {
   "value": %q
  },
  "tags": {
   "value": {}
  }
}
PARAM
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
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

func (SubscriptionTemplateDeploymentResource) singleItemWithParameterConfigAndVariable(data acceptance.TestData, templateVar string, paramVal string) string {
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
