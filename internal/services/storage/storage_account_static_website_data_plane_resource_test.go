// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AccountStaticWebsiteResource struct{}

func TestAccountStaticWebsiteResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_static_website", "test")
	r := AccountStaticWebsiteResource{}

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

func TestAccountStaticWebsiteResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_static_website", "test")
	r := AccountStaticWebsiteResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withIndex(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.with404(data),
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

func TestAccountStaticWebsiteResource_with404(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_static_website", "test")
	r := AccountStaticWebsiteResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.with404(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccountStaticWebsiteResource_withIndex(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_static_website", "test")
	r := AccountStaticWebsiteResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIndex(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AccountStaticWebsiteResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseStorageAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	accountDetails, err := client.Storage.GetAccount(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if accountDetails == nil {
		return nil, fmt.Errorf("unable to locate %s", *id)
	}

	accountsClient, err := client.Storage.AccountsDataPlaneClient(ctx, *accountDetails, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Accounts Data Plane Client for %s: %+v", *id, err)
	}

	props, err := accountsClient.GetServiceProperties(ctx, id.StorageAccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving static website properties for %s: %+v", *id, err)
	}

	if props.StaticWebsite == nil {
		return nil, nil
	}

	return pointer.To(props.StaticWebsite.Enabled), nil
}

func (r AccountStaticWebsiteResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_static_website" "test" {
  storage_account_id = azurerm_storage_account.test.id
  error_404_document = "sadpanda_2.html"
  index_document     = "index_2.html"
}
`, r.template(data))
}

func (r AccountStaticWebsiteResource) with404(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_static_website" "test" {
  storage_account_id = azurerm_storage_account.test.id
  error_404_document = "sadpanda.html"
}
`, r.template(data))
}

func (r AccountStaticWebsiteResource) withIndex(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_static_website" "test" {
  storage_account_id = azurerm_storage_account.test.id
  index_document     = "index.html"
}
`, r.template(data))
}

func (r AccountStaticWebsiteResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}


resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
