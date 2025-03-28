// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DevCenterCatalogsResource struct{}

func (r DevCenterCatalogsResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := catalogs.ParseDevCenterCatalogID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.DevCenter.V20250201.Catalogs.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func TestAccDevCenterCatalogs_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_catalog", "test")
	r := DevCenterCatalogsResource{}

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

func TestAccDevCenterCatalogs_adoGit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_catalog", "test")
	r := DevCenterCatalogsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.adoGit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDevCenterCatalogs_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_catalog", "test")
	r := DevCenterCatalogsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func (r DevCenterCatalogsResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_dev_center_catalog" "test" {
  name                = "acctest-catalog-%d"
  resource_group_name = azurerm_resource_group.test.name
  dev_center_id       = azurerm_dev_center.test.id
  catalog_github {
    branch            = "main"
    path              = "/template"
    uri               = "https://github.com/am-lim/deployment-environments.git"
    key_vault_key_url = "https://amlim-kv.vault.azure.net/secrets/envTest/0a79f15246ce4b35a13957367b422cab"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterCatalogsResource) adoGit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_dev_center_catalog" "test" {
  name                = "acctest-catalog-%d"
  resource_group_name = azurerm_resource_group.test.name
  dev_center_id       = azurerm_dev_center.test.id
  catalog_adogit {
    branch            = "main"
    path              = "/template"
    uri               = "https://amlim@dev.azure.com/amlim/testCatalog/_git/testCatalog"
    key_vault_key_url = "https://amlim-kv.vault.azure.net/secrets/ado/6279752c2bdd4a38a3e79d958cc36a75"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterCatalogsResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_dev_center_catalog" "test" {
  name                = "acctest-catalog-%d"
  resource_group_name = azurerm_resource_group.test.name
  dev_center_id       = azurerm_dev_center.test.id
  catalog_github {
    branch            = "foo"
    path              = ""
    uri               = "https://github.com/am-lim/deployment-environments.git"
    key_vault_key_url = "https://amlim-kv.vault.azure.net/secrets/envTest/0a79f15246ce4b35a13957367b422cab"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DevCenterCatalogsResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%s"
}

resource "azurerm_dev_center" "test" {
  name                = "acctdc-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, "West Europe")
}
