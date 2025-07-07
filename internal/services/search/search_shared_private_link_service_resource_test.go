// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package search_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2024-06-01-preview/sharedprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SearchSharedPrivateLinkServiceResource struct{}

func TestAccSearchSharedPrivateLinkServiceResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_shared_private_link_service", "test")
	r := SearchSharedPrivateLinkServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("request_message").HasValue("please approve")),
		},
		data.ImportStep(),
	})
}

func TestAccSearchSharedPrivateLinkServiceResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_shared_private_link_service", "test")
	r := SearchSharedPrivateLinkServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r SearchSharedPrivateLinkServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sharedprivatelinkresources.ParseSharedPrivateLinkResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Search.SearchSharedPrivateLinkResourceClient.Get(ctx, *id, sharedprivatelinkresources.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s was not found: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r SearchSharedPrivateLinkServiceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_search_shared_private_link_service" "test" {
  name               = "acctest%[2]d"
  search_service_id  = azurerm_search_service.test.id
  subresource_name   = "blob"
  target_resource_id = azurerm_storage_account.test.id
  request_message    = "please approve"
}

resource "azurerm_search_shared_private_link_service" "test2" {
  name               = "acctest2%[2]d"
  search_service_id  = azurerm_search_service.test.id
  subresource_name   = "file"
  target_resource_id = azurerm_storage_account.test.id
  request_message    = "please approve"
}
`, template, data.RandomInteger)
}

func (r SearchSharedPrivateLinkServiceResource) requiresImport(data acceptance.TestData) string {
	template := SearchSharedPrivateLinkServiceResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_search_shared_private_link_service" "import" {
  name               = azurerm_search_shared_private_link_service.test.name
  search_service_id  = azurerm_search_shared_private_link_service.test.search_service_id
  subresource_name   = azurerm_search_shared_private_link_service.test.subresource_name
  target_resource_id = azurerm_search_shared_private_link_service.test.target_resource_id
  request_message    = azurerm_search_shared_private_link_service.test.request_message
}
`, template)
}

func (r SearchSharedPrivateLinkServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_search_service" "test" {
  name                = "acctestsearchservice%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account" "test" {
  name                = "acc%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
