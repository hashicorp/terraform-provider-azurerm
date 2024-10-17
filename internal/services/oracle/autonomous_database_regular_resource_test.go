// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AdbsRegularResource struct{}

func (a AdbsRegularResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Oracle.OracleClient.AutonomousDatabases.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving adbs %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAdbsRegularResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseRegularResource{}.ResourceType(), "test")
	r := AdbsRegularResource{}
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

func TestAdbsRegularResource_update(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseRegularResource{}.ResourceType(), "test")
	r := AdbsRegularResource{}
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

func TestAdbsRegularResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.AutonomousDatabaseRegularResource{}.ResourceType(), "test")
	r := AdbsRegularResource{}
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

func (a AdbsRegularResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_oracle_autonomous_database" "test" {
  name = "OFake%[2]d"

  display_name = "OFake%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location = "%[3]s"
  compute_model = "ECPU"
  compute_count = "2"
  license_model = "BringYourOwnLicense"
  backup_retention_period_in_days = 12
  is_auto_scaling_enabled = false
  is_auto_scaling_for_storage_enabled = false
  is_mtls_connection_required = false
  data_storage_size_in_gbs = "32"
  db_workload = "OLTP"
  admin_password = "TestPass#2024#"
  db_version = "19c"
  character_set = "AL32UTF8"
  ncharacter_set = "AL16UTF16"
  subnet_id = azurerm_subnet.virtual_network_subnet.id
  vnet_id = azurerm_virtual_network.virtual_network.id
  lifecycle {
    ignore_changes = [
      admin_password
    ]
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a AdbsRegularResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_oracle_autonomous_database" "test" {
  name = "OFake%[2]d"
  display_name = "OFake%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location = "%[3]s"
  compute_model = "ECPU"
  compute_count = "2"
  license_model = "BringYourOwnLicense"
  backup_retention_period_in_days = 12
  is_auto_scaling_enabled = false
  is_auto_scaling_for_storage_enabled = false
  is_mtls_connection_required = false
  data_storage_size_in_gbs = "32"
  db_workload = "OLTP"
  admin_password = "TestPass#2024#"
  db_version = "19c"
  character_set = "AL32UTF8"
  ncharacter_set = "AL16UTF16"
  subnet_id = azurerm_subnet.virtual_network_subnet.id
  vnet_id = azurerm_virtual_network.virtual_network.id
  tags = {
    test = "test1"
  }
  lifecycle {
    ignore_changes = [
      admin_password
    ]
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a AdbsRegularResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_oracle_autonomous_database" "import" {
  name = azurerm_oracle_autonomous_database.test.name
  display_name = azurerm_oracle_autonomous_database.test.display_name
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location = azurerm_oracle_autonomous_database.test.location
  compute_model = azurerm_oracle_autonomous_database.test.compute_model
  compute_count = azurerm_oracle_autonomous_database.test.compute_count
  license_model = azurerm_oracle_autonomous_database.test.license_model
  backup_retention_period_in_days = azurerm_oracle_autonomous_database.test.backup_retention_period_in_days
  is_auto_scaling_enabled = azurerm_oracle_autonomous_database.test.is_auto_scaling_enabled
  is_auto_scaling_for_storage_enabled = azurerm_oracle_autonomous_database.test.is_auto_scaling_for_storage_enabled
  is_mtls_connection_required = azurerm_oracle_autonomous_database.test.is_mtls_connection_required
  data_storage_size_in_gbs = azurerm_oracle_autonomous_database.test.data_storage_size_in_gbs
  db_workload = azurerm_oracle_autonomous_database.test.db_workload
  admin_password = azurerm_oracle_autonomous_database.test.admin_password
  db_version = azurerm_oracle_autonomous_database.test.db_version
  character_set = azurerm_oracle_autonomous_database.test.character_set
  ncharacter_set = azurerm_oracle_autonomous_database.test.ncharacter_set
  subnet_id = azurerm_oracle_autonomous_database.test.subnet_id
  vnet_id = azurerm_oracle_autonomous_database.test.vnet_id
  lifecycle {
    ignore_changes = [
      admin_password
    ]
  }
}
`, a.basic(data))
}

func (a AdbsRegularResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "virtual_network" {
  name                = "OFakeacctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "virtual_network_subnet" {
  name                 = "OFakeacctest%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.virtual_network.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      actions = [
        "Microsoft.Network/networkinterfaces/*",
		"Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Oracle.Database/networkAttachments"
    }
  }
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
