// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/groupuser"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementGroupUserResource struct{}

func TestAccAzureRMApiManagementGroupUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group_user", "test")
	r := ApiManagementGroupUserResource{}

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

func TestAccAzureRMApiManagementGroupUser_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group_user", "test")
	r := ApiManagementGroupUserResource{}

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

func (ApiManagementGroupUserResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := groupuser.ParseGroupUserID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.GroupUsersClient.CheckEntityExists(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	// the HEAD API not found returns resp 404, but no err
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}

	return pointer.To(true), nil
}

func (ApiManagementGroupUserResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

  sku_name = "Developer_1"
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Test Group"
}

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
}

resource "azurerm_api_management_group_user" "test" {
  user_id             = azurerm_api_management_user.test.user_id
  group_name          = azurerm_api_management_group.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementGroupUserResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_group_user" "import" {
  user_id             = azurerm_api_management_group_user.test.user_id
  group_name          = azurerm_api_management_group_user.test.group_name
  api_management_name = azurerm_api_management_group_user.test.api_management_name
  resource_group_name = azurerm_api_management_group_user.test.resource_group_name
}
`, r.basic(data))
}
