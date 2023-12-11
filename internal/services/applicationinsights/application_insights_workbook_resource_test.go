// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	workbooks "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApplicationInsightsWorkbookResource struct{}

func TestAccApplicationInsightsWorkbook_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")
	r := ApplicationInsightsWorkbookResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, data.RandomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsightsWorkbook_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")
	r := ApplicationInsightsWorkbookResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, data.RandomInteger),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApplicationInsightsWorkbook_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")
	r := ApplicationInsightsWorkbookResource{}
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

func TestAccApplicationInsightsWorkbook_updateDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")
	r := ApplicationInsightsWorkbookResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationInsightsWorkbook_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")
	r := ApplicationInsightsWorkbookResource{}
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

func TestAccApplicationInsightsWorkbook_hiddenTitleInTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")
	r := ApplicationInsightsWorkbookResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.hiddenTitleInTags(data),
			ExpectError: regexp.MustCompile("a tag with the key `hidden-title` should not be used to set the display name. Please Use `display_name` instead"),
		},
	})
}

func (r ApplicationInsightsWorkbookResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workbooks.ParseWorkbookID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.AppInsights.WorkbookClient
	resp, err := client.WorkbooksGet(ctx, *id, workbooks.WorkbooksGetOperationOptions{CanFetchContent: utils.Bool(true)})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ApplicationInsightsWorkbookResource) template(data acceptance.TestData) string {
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

func (r ApplicationInsightsWorkbookResource) basic(data acceptance.TestData, intValue int) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_application_insights_workbook" "test" {
  name                = "be1ad266-d329-4454-b693-8287e4d3b35d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "acctest-amw-%d"
  data_json = jsonencode({
    "version" = "Notebook/1.0",
    "items" = [
      {
        "type" = 1,
        "content" = {
          "json" = "Test2022"
        },
        "name" = "text - 0"
      }
    ],
    "isLocked" = false,
    "fallbackResourceIds" = [
      "Azure Monitor"
    ]
  })
}
`, template, intValue)
}

func (r ApplicationInsightsWorkbookResource) hiddenTitleInTags(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_workbook" "test" {
  name                = "be1ad266-d329-4454-b693-8287e4d3b35d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "acctest-amw-%d"
  data_json = jsonencode({
    "version" = "Notebook/1.0",
    "items" = [
      {
        "type" = 1,
        "content" = {
          "json" = "Test2022"
        },
        "name" = "text - 0"
      }
    ],
    "isLocked" = false,
    "fallbackResourceIds" = [
      "Azure Monitor"
    ]
  })
  tags = {
    hidden-title = "Test Display Name"
  }
}
`, template, data.RandomInteger)
}

func (r ApplicationInsightsWorkbookResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data, data.RandomInteger)
	return fmt.Sprintf(`
			%s

resource "azurerm_application_insights_workbook" "import" {
  name                = azurerm_application_insights_workbook.test.name
  resource_group_name = azurerm_application_insights_workbook.test.resource_group_name
  location            = azurerm_application_insights_workbook.test.location
  category            = azurerm_application_insights_workbook.test.category
  display_name        = azurerm_application_insights_workbook.test.display_name
  source_id           = azurerm_application_insights_workbook.test.source_id
  data_json           = azurerm_application_insights_workbook.test.data_json
}
`, config)
}

func (r ApplicationInsightsWorkbookResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsads%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_application_insights_workbook" "test" {
  name                 = "0f498fab-2989-4395-b084-fc092d83a6b1"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  display_name         = "acctest-amw-1"
  source_id            = lower(azurerm_resource_group.test.id)
  category             = "workbook1"
  description          = "description1"
  storage_container_id = azurerm_storage_container.test.resource_manager_id

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  data_json = jsonencode({
    "version" = "Notebook/1.0",
    "items" = [
      {
        "type" = 1,
        "content" = {
          "json" = "Test2021"
        },
        "name" = "text - 0"
      }
    ],
    "isLocked" = false,
    "fallbackResourceIds" = [
      "Azure Monitor"
    ]
  })
  tags = {
    env = "test"
  }

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, template, data.RandomInteger, data.RandomString)
}

func (r ApplicationInsightsWorkbookResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsads%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_application_insights_workbook" "test" {
  name                 = "0f498fab-2989-4395-b084-fc092d83a6b1"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  display_name         = "acctest-amw-2"
  source_id            = "azure monitor"
  category             = "workbook2"
  description          = "description2"
  storage_container_id = azurerm_storage_container.test.resource_manager_id

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  data_json = jsonencode({
    "version" = "Notebook/1.0",
    "items" = [
      {
        "type" = 1,
        "content" = {
          "json" = "Test2022"
        },
        "name" = "text - 0"
      }
    ],
    "isLocked" = false,
    "fallbackResourceIds" = [
      "Azure Monitor"
    ]
  })
  tags = {
    env = "test2"
  }

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, template, data.RandomInteger, data.RandomString)
}
