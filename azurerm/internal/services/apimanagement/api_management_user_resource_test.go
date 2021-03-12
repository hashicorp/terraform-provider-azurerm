package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementUserResource struct {
}

func TestAccApiManagementUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")
	r := ApiManagementUserResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementUser_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")
	r := ApiManagementUserResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("first_name").HasValue("Acceptance"),
				check.That(data.ResourceName).Key("last_name").HasValue("Test"),
				check.That(data.ResourceName).Key("state").HasValue("active"),
			),
		},
		{
			Config: r.updatedBlocked(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("first_name").HasValue("Acceptance Updated"),
				check.That(data.ResourceName).Key("last_name").HasValue("Test Updated"),
				check.That(data.ResourceName).Key("state").HasValue("blocked"),
			),
		},
		{
			Config: r.updatedActive(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.password(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.invited(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.signUp(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
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

func (ApiManagementUserResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	userId := id.Path["users"]

	resp, err := clients.ApiManagement.UsersClient.Get(ctx, resourceGroup, serviceName, userId)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement User (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
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
