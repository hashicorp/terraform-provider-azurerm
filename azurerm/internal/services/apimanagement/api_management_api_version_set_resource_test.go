package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementApiVersionSetResource struct {
}

func TestAccApiManagementApiVersionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiVersionSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

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

func TestAccApiManagementApiVersionSet_header(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.header(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiVersionSet_query(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.query(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiVersionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("TestDescription1"),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("TestApiVersionSet1%d", data.RandomInteger)),
			),
		},
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("TestDescription2"),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("TestApiVersionSet2%d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func (t ApiManagementApiVersionSetResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ApiVersionSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiVersionSetClient.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagementApi Version Set (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ApiManagementApiVersionSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_scheme   = "Segment"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiVersionSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "import" {
  name                = azurerm_api_management_api_version_set.test.name
  resource_group_name = azurerm_api_management_api_version_set.test.resource_group_name
  api_management_name = azurerm_api_management_api_version_set.test.api_management_name
  description         = azurerm_api_management_api_version_set.test.description
  display_name        = azurerm_api_management_api_version_set.test.display_name
  versioning_scheme   = azurerm_api_management_api_version_set.test.versioning_scheme
}
`, r.basic(data))
}

func (r ApiManagementApiVersionSetResource) header(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_scheme   = "Header"
  version_header_name = "Header1"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiVersionSetResource) query(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_scheme   = "Query"
  version_query_name  = "Query1"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiVersionSetResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  description         = "TestDescription2"
  display_name        = "TestApiVersionSet2%d"
  versioning_scheme   = "Segment"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (ApiManagementApiVersionSetResource) template(data acceptance.TestData) string {
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
  sku_name            = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
