// Copyright © 2025, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exascaledbstoragevaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ExascaleDbStorageVaultResource struct{}

func (a ExascaleDbStorageVaultResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := exascaledbstoragevaults.ParseExascaleDbStorageVaultID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient25.ExascaleDbStorageVaults.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving db storage vault %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestDbStorageVaultResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExascaleDbStorageVaultResource{}.ResourceType(), "test")
	r := ExascaleDbStorageVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("high_capacity_database_storage_input", "high_capacity_database_storage"),
	})
}

func TestDbStorageVaultResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExascaleDbStorageVaultResource{}.ResourceType(), "test")
	r := ExascaleDbStorageVaultResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("high_capacity_database_storage_input", "high_capacity_database_storage"),
	})
}

func TestDbStorageVaultResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.ExascaleDbStorageVaultResource{}.ResourceType(), "test")
	r := ExascaleDbStorageVaultResource{}
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

func (a ExascaleDbStorageVaultResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_exascale_db_storage_vault" "test" {
  name                             = "OFake%[2]d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = "%[3]s"
  zones               			   = ["2"]
  display_name                     = "OFake%[2]d"
  description                      = "description"
  additional_flash_cache_in_percent = 100
  high_capacity_database_storage_input {
    total_size_in_gbs = 300
  }
  time_zone        			  = "UTC"
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExascaleDbStorageVaultResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_exascale_db_storage_vault" "test" {
  name                             = "OFake%[2]d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = "%[3]s"
  zones                            = ["2"]
  display_name                     = "OFake%[2]d"
  high_capacity_database_storage_input {
	total_size_in_gbs = 300
  }
  additional_flash_cache_in_percent = 100
  description                    = "description"
  high_capacity_database_storage {
	total_size_in_gbs = 300
	available_size_in_gbs = 60
  }
  time_zone        			  = "UTC"
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a ExascaleDbStorageVaultResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_exascale_db_storage_vault" "import" {
  name                          	   = azurerm_oracle_exascale_db_storage_vault.test.name
  resource_group_name            	   = azurerm_oracle_exascale_db_storage_vault.test.resource_group_name
  location                             = azurerm_oracle_exascale_db_storage_vault.test.location
  display_name                         = azurerm_oracle_exascale_db_storage_vault.test.display_name
  additional_flash_cache_in_percent	   = azurerm_oracle_exascale_db_storage_vault.test.additional_flash_cache_in_percent
  description                      	   = azurerm_oracle_exascale_db_storage_vault.test.description
  time_zone  					       = azurerm_oracle_exascale_db_storage_vault.test.time_zone
  zones                                = azurerm_oracle_exascale_db_storage_vault.test.zones
}
`, a.basic(data))
}

func (a ExascaleDbStorageVaultResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
