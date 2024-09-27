// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementUserResource struct{}

func TestAccApiManagementUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")
	r := ApiManagementUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("first_name").HasValue("Acceptance"),
				check.That(data.ResourceName).Key("last_name").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementUser_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")
	r := ApiManagementUserResource{}

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

func TestAccApiManagementUser_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")
	r := ApiManagementUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("first_name").HasValue("Acceptance"),
				check.That(data.ResourceName).Key("last_name").HasValue("Test"),
				check.That(data.ResourceName).Key("state").HasValue("active"),
			),
		},
		{
			Config: r.updatedBlocked(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("first_name").HasValue("Acceptance Updated"),
				check.That(data.ResourceName).Key("last_name").HasValue("Test Updated"),
				check.That(data.ResourceName).Key("state").HasValue("blocked"),
			),
		},
		{
			Config: r.updatedActive(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("first_name").HasValue("Acceptance"),
				check.That(data.ResourceName).Key("last_name").HasValue("Test"),
				check.That(data.ResourceName).Key("state").HasValue("active"),
			),
		},
	})
}

func TestAccApiManagementUser_password(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")
	r := ApiManagementUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.password(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("first_name").HasValue("Acceptance"),
				check.That(data.ResourceName).Key("last_name").HasValue("Test"),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"password"},
		},
	})
}

func TestAccApiManagementUser_invite(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")
	r := ApiManagementUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.invited(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				// not returned
				"confirmation",
			},
		},
	})
}

func TestAccApiManagementUser_signup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")
	r := ApiManagementUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.signUp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				// not returned
				"confirmation",
			},
		},
	})
}

func TestAccApiManagementUser_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")
	r := ApiManagementUserResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("first_name").HasValue("Acceptance"),
				check.That(data.ResourceName).Key("last_name").HasValue("Test"),
				check.That(data.ResourceName).Key("note").HasValue("Used for testing in dimension C-137."),
			),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				// not returned
				"confirmation",
			},
		},
	})
}

func (ApiManagementUserResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := user.ParseUserID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.UsersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementUserResource) basic(data acceptance.TestData) string {
	template := ApiManagementUserResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementUserResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "import" {
  user_id             = azurerm_api_management_user.test.user_id
  api_management_name = azurerm_api_management_user.test.api_management_name
  resource_group_name = azurerm_api_management_user.test.resource_group_name
  first_name          = azurerm_api_management_user.test.first_name
  last_name           = azurerm_api_management_user.test.last_name
  email               = azurerm_api_management_user.test.email
  state               = azurerm_api_management_user.test.state
}
`, r.basic(data))
}

func (r ApiManagementUserResource) password(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
  password            = "3991bb15-282d-4b9b-9de3-3d5fc89eb530"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementUserResource) updatedActive(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementUserResource) updatedBlocked(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance Updated"
  last_name           = "Test Updated"
  email               = "azure-acctest%d@example.com"
  state               = "blocked"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementUserResource) invited(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test User"
  email               = "azure-acctest%d@example.com"
  state               = "blocked"
  confirmation        = "invite"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementUserResource) signUp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test User"
  email               = "azure-acctest%d@example.com"
  state               = "blocked"
  confirmation        = "signup"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementUserResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
  confirmation        = "signup"
  note                = "Used for testing in dimension C-137."
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (ApiManagementUserResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
