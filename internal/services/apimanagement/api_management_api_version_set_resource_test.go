// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apiversionset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiVersionSetResource struct{}

func TestAccApiManagementApiVersionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

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

func TestAccApiManagementApiVersionSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

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

func TestAccApiManagementApiVersionSet_header(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.header(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiVersionSet_query(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.query(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiVersionSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_version_set", "test")
	r := ApiManagementApiVersionSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("TestDescription1"),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("TestApiVersionSet1%d", data.RandomInteger)),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("TestDescription2"),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("TestApiVersionSet2%d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementApiVersionSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apiversionset.ParseApiVersionSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiVersionSetClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
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
  sku_name            = "Consumption_0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
