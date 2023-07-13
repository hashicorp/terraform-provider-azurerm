// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	workbooktemplates "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApplicationInsightsWorkbookTemplateResource struct{}

func TestAccApplicationInsightsWorkbookTemplate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook_template", "test")
	r := ApplicationInsightsWorkbookTemplateResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsightsWorkbookTemplate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook_template", "test")
	r := ApplicationInsightsWorkbookTemplateResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApplicationInsightsWorkbookTemplate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook_template", "test")
	r := ApplicationInsightsWorkbookTemplateResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsightsWorkbookTemplate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook_template", "test")
	r := ApplicationInsightsWorkbookTemplateResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ApplicationInsightsWorkbookTemplateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workbooktemplates.ParseWorkbookTemplateID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.AppInsights.WorkbookTemplateClient
	resp, err := client.WorkbookTemplatesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ApplicationInsightsWorkbookTemplateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApplicationInsightsWorkbookTemplateResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_application_insights_workbook_template" "test" {
  name                = "acctest-aiwt-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"

  galleries {
    category = "workbook"
    name     = "test"
  }

  template_data = jsonencode({
    "version" : "Notebook/1.0",
    "items" : [
      {
        "type" : 1,
        "content" : {
          "json" : "## New workbook\n---\n\nWelcome to your new workbook."
        },
        "name" : "text - 2"
      }
    ],
    "styleSettings" : {},
    "$schema" : "https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json"
  })
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r ApplicationInsightsWorkbookTemplateResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_application_insights_workbook_template" "import" {
  name                = azurerm_application_insights_workbook_template.test.name
  resource_group_name = azurerm_application_insights_workbook_template.test.resource_group_name
  location            = azurerm_application_insights_workbook_template.test.location

  galleries {
    category = "workbook"
    name     = "test"
  }

  template_data = jsonencode({
    "version" : "Notebook/1.0",
    "items" : [
      {
        "type" : 1,
        "content" : {
          "json" : "## New workbook\n---\n\nWelcome to your new workbook."
        },
        "name" : "text - 2"
      }
    ],
    "styleSettings" : {},
    "$schema" : "https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json"
  })
}
`, config)
}

func (r ApplicationInsightsWorkbookTemplateResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_application_insights_workbook_template" "test" {
  name                = "acctest-aiwt-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  author              = "test author"
  priority            = 1

  galleries {
    category      = "Failures"
    name          = "test"
    order         = 100
    resource_type = "microsoft.insights/components"
    type          = "tsg"
  }

  template_data = jsonencode({
    "version" : "Notebook/1.0",
    "items" : [
      {
        "type" : 1,
        "content" : {
          "json" : "## New workbook\n---\n\nWelcome to your new workbook."
        },
        "name" : "text - 2"
      }
    ],
    "styleSettings" : {},
    "$schema" : "https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json"
  })

  localized = jsonencode({
    "ar" : [
      {
        "galleries" : [
          {
            "name" : "test",
            "category" : "Failures",
            "type" : "tsg",
            "resourceType" : "microsoft.insights/components",
            "order" : 100
          }
        ],
        "templateData" : {
          "version" : "Notebook/1.0",
          "items" : [
            {
              "type" : 1,
              "content" : {
                "json" : "## New workbook\n---\n\nWelcome to your new workbook."
              },
              "name" : "text - 2"
            }
          ],
          "styleSettings" : {},
          "$schema" : "https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json"
        },
      }
    ]
  })

  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r ApplicationInsightsWorkbookTemplateResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_application_insights_workbook_template" "test" {
  name                = "acctest-aiwt-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  author              = "test author 2"
  priority            = 2

  galleries {
    category      = "workbook"
    name          = "test2"
    order         = 200
    resource_type = "Azure Monitor"
    type          = "workbook"
  }

  template_data = jsonencode({
    "version" : "Notebook/1.0",
    "items" : [
      {
        "type" : 2,
        "content" : {
          "json" : "## New workbook\n---\n\nWelcome to your new workbook."
        },
        "name" : "text - 2"
      }
    ],
    "styleSettings" : {},
    "$schema" : "https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json"
  })

  localized = jsonencode({
    "en-US" : [
      {
        "galleries" : [
          {
            "name" : "test2",
            "category" : "workbook",
            "type" : "workbook",
            "resourceType" : "Azure Monitor",
            "order" : 200
          }
        ],
        "templateData" : {
          "version" : "Notebook/1.0",
          "items" : [
            {
              "type" : 2,
              "content" : {
                "json" : "## New workbook\n---\n\nWelcome to your new workbook."
              },
              "name" : "text - 2"
            }
          ],
          "styleSettings" : {},
          "$schema" : "https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json"
        },
      }
    ]
  })

  tags = {
    key = "value2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
