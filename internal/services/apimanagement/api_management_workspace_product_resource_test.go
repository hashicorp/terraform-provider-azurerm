// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/product"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceProductTestResource struct{}

func TestAccApiManagementWorkspaceProduct_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_product", "test")
	r := ApiManagementWorkspaceProductTestResource{}

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

func TestAccApiManagementWorkspaceProduct_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_product", "test")
	r := ApiManagementWorkspaceProductTestResource{}

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

func TestAccApiManagementWorkspaceProduct_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_product", "test")
	r := ApiManagementWorkspaceProductTestResource{}

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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspaceProduct_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_product", "test")
	r := ApiManagementWorkspaceProductTestResource{}

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

func TestAccApiManagementWorkspaceProduct_specifySubscriptionsLimitError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_product", "test")
	r := ApiManagementWorkspaceProductTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.specifySubscriptionsLimitError(data),
			ExpectError: regexp.MustCompile("`require_subscription_enabled` must be set to `true` when `subscriptions_limit` is specified"),
		},
	})
}

func TestAccApiManagementWorkspaceProduct_specifyRequireApprovalEnabledError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_product", "test")
	r := ApiManagementWorkspaceProductTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.specifyRequireApprovalEnabledError(data),
			ExpectError: regexp.MustCompile("`require_subscription_enabled` must be set to `true` when `require_approval_enabled` is specified"),
		},
	})
}

func (ApiManagementWorkspaceProductTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := product.ParseWorkspaceProductID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ProductClient_v2024_05_01.WorkspaceProductGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceProductTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_product" "test" {
  name                        = "acctestAMProduct-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Test Product"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceProductTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_product" "import" {
  name                        = azurerm_api_management_workspace_product.test.name
  api_management_workspace_id = azurerm_api_management_workspace_product.test.api_management_workspace_id
  display_name                = azurerm_api_management_workspace_product.test.display_name
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceProductTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_product" "test" {
  name                         = "acctestAMProduct-%d"
  api_management_workspace_id  = azurerm_api_management_workspace.test.id
  display_name                 = "Test Product Complete"
  published_enabled            = true
  description                  = "Test Product Complete Description"
  terms                        = "Test Product Complete Terms"
  require_subscription_enabled = true
  require_approval_enabled     = true
  subscriptions_limit          = 5
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceProductTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_product" "test" {
  name                         = "acctestAMProduct-%d"
  api_management_workspace_id  = azurerm_api_management_workspace.test.id
  display_name                 = "Test Product Updated"
  published_enabled            = false
  description                  = "Test Product Description"
  terms                        = "Test Product Terms"
  require_subscription_enabled = true
  require_approval_enabled     = false
  subscriptions_limit          = 6
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceProductTestResource) specifySubscriptionsLimitError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_product" "test" {
  name                         = "acctestAMProduct-%d"
  api_management_workspace_id  = azurerm_api_management_workspace.test.id
  display_name                 = "Test Product"
  published_enabled            = true
  description                  = "Test Product Description"
  terms                        = "Test Product Terms"
  require_subscription_enabled = false
  subscriptions_limit          = 5
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceProductTestResource) specifyRequireApprovalEnabledError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_product" "test" {
  name                         = "acctestAMProduct-%d"
  api_management_workspace_id  = azurerm_api_management_workspace.test.id
  display_name                 = "Test Product"
  published_enabled            = true
  description                  = "Test Product Description"
  terms                        = "Test Product Terms"
  require_subscription_enabled = false
  require_approval_enabled     = true
}
`, r.template(data), data.RandomInteger)
}

func (ApiManagementWorkspaceProductTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestAMWorkspace-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
