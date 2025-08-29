// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspaceApiVersionSetResource struct{}

func TestAccApiManagementWorkspaceApiVersionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api_version_set", "test")
	r := ApiManagementWorkspaceApiVersionSetResource{}

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

func TestAccApiManagementWorkspaceApiVersionSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api_version_set", "test")
	r := ApiManagementWorkspaceApiVersionSetResource{}

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

func TestAccApiManagementWorkspaceApiVersionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api_version_set", "test")
	r := ApiManagementWorkspaceApiVersionSetResource{}

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

func TestAccApiManagementWorkspaceApiVersionSet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_api_version_set", "test")
	r := ApiManagementWorkspaceApiVersionSetResource{}

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

func (ApiManagementWorkspaceApiVersionSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apiversionset.ParseWorkspaceApiVersionSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiVersionSetClient_v2024_05_01.WorkspaceApiVersionSetGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspaceApiVersionSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api_version_set" "test" {
  name                        = "acctestAMWAVS-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Test API Version Set"
  versioning_scheme           = "Segment"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiVersionSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace_api_version_set" "import" {
  name                        = azurerm_api_management_workspace_api_version_set.test.name
  api_management_workspace_id = azurerm_api_management_workspace_api_version_set.test.api_management_workspace_id
  display_name                = azurerm_api_management_workspace_api_version_set.test.display_name
  versioning_scheme           = azurerm_api_management_workspace_api_version_set.test.versioning_scheme
}
`, r.basic(data))
}

func (r ApiManagementWorkspaceApiVersionSetResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api_version_set" "test" {
  name                        = "acctestAMWAVS-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Updated API Version Set"
  versioning_scheme           = "Header"
  description                 = "Updated description"
  version_header_name         = "Api-Version"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementWorkspaceApiVersionSetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_api_version_set" "test" {
  name                        = "acctestAMWAVS-%d"
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  display_name                = "Complete API Version Set"
  versioning_scheme           = "Query"
  description                 = "A complete API version set"
  version_query_name          = "version"
}
`, r.template(data), data.RandomInteger)
}

func (ApiManagementWorkspaceApiVersionSetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
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
  name              = "acctestAMW-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "Test workspace description"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
