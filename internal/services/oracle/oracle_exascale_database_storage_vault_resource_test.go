// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exascaledbstoragevaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExascaleDatabaseStorageVaultResource struct{}

func (a ExascaleDatabaseStorageVaultResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := exascaledbstoragevaults.ParseExascaleDbStorageVaultID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient.ExascaleDbStorageVaults.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestDbStorageVaultResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExascaleDatabaseStorageVaultResource{}.ResourceType(), "test")
	r := ExascaleDatabaseStorageVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("high_capacity_database_storage.#").HasValue("1"),
			),
		},
	})
}

func TestDbStorageVaultResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExascaleDatabaseStorageVaultResource{}.ResourceType(), "test")
	r := ExascaleDatabaseStorageVaultResource{}
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

func TestDbStorageVaultResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExascaleDatabaseStorageVaultResource{}.ResourceType(), "test")
	r := ExascaleDatabaseStorageVaultResource{}
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestDbStorageVaultResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExascaleDatabaseStorageVaultResource{}.ResourceType(), "test")
	r := ExascaleDatabaseStorageVaultResource{}
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

func (a ExascaleDatabaseStorageVaultResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_exascale_database_storage_vault" "test" {
  name                              = "OFakeacctest%[2]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "%[3]s"
  zones                             = ["2"]
  display_name                      = "OFakeacctest%[2]d"
  additional_flash_cache_percentage = 100
  high_capacity_database_storage {
    total_size_in_gb = 300
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExascaleDatabaseStorageVaultResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_exascale_database_storage_vault" "test" {
  name                = "OFakeacctest%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  zones               = ["2"]
  display_name        = "OFakeacctest%[2]d"
  high_capacity_database_storage {
    total_size_in_gb = 300
  }
  additional_flash_cache_percentage = 100
  description                       = "description"
  time_zone                         = "America/New_York"
  tags = {
    ENV = "Test"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExascaleDatabaseStorageVaultResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_exascale_database_storage_vault" "test" {
  name                              = "OFakeacctest%[2]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "%[3]s"
  zones                             = ["2"]
  display_name                      = "OFakeacctest%[2]d"
  additional_flash_cache_percentage = 100
  high_capacity_database_storage {
    total_size_in_gb = 300
  }
  tags = {
    ENV = "Test"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExascaleDatabaseStorageVaultResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_exascale_database_storage_vault" "import" {
  name                              = azurerm_oracle_exascale_database_storage_vault.test.name
  resource_group_name               = azurerm_oracle_exascale_database_storage_vault.test.resource_group_name
  location                          = azurerm_oracle_exascale_database_storage_vault.test.location
  display_name                      = azurerm_oracle_exascale_database_storage_vault.test.display_name
  additional_flash_cache_percentage = azurerm_oracle_exascale_database_storage_vault.test.additional_flash_cache_percentage
  description                       = azurerm_oracle_exascale_database_storage_vault.test.description
  time_zone                         = azurerm_oracle_exascale_database_storage_vault.test.time_zone
  zones                             = azurerm_oracle_exascale_database_storage_vault.test.zones
  high_capacity_database_storage {
    total_size_in_gb = azurerm_oracle_exascale_database_storage_vault.test.high_capacity_database_storage[0].total_size_in_gb
  }
}
`, a.basic(data))
}

func (a ExascaleDatabaseStorageVaultResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary)
}
