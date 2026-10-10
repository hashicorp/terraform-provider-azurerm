// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage"
)

type StorageAccountSasDataSource struct{}

func TestAccDataSourceStorageAccountSas_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account_sas", "test")
	ipAddresses := "10.0.0.1-10.0.0.4"
	utcNow := time.Now().UTC()
	startDate := utcNow.Format(time.RFC3339)
	endDate := utcNow.Add(time.Hour * 24).Format(time.RFC3339)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountSasDataSource{}.basic(data, startDate, endDate, ipAddresses),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("https_only").HasValue("true"),
				check.That(data.ResourceName).Key("ip_addresses").HasValue(ipAddresses),
				check.That(data.ResourceName).Key("signed_version").HasValue("2019-10-10"),
				check.That(data.ResourceName).Key("start").HasValue(startDate),
				check.That(data.ResourceName).Key("expiry").HasValue(endDate),
				check.That(data.ResourceName).Key("sas").Exists(),
			),
		},
	})
}

func TestAccDataSourceStorageAccountSas_noPermissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account_sas", "test")
	ipAddresses := "10.0.0.1-10.0.0.4"
	utcNow := time.Now().UTC()
	startDate := utcNow.Format(time.RFC3339)
	endDate := utcNow.Add(time.Hour * 24).Format(time.RFC3339)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountSasDataSource{}.noPermissions(data, startDate, endDate, ipAddresses),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sas").Exists(),
			),
		},
	})
}

func (d StorageAccountSasDataSource) basic(data acceptance.TestData, startDate string, endDate string, ipAddresses string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true
  ip_addresses      = "%s"
  signed_version    = "2019-10-10"

  resource_types {
    service   = true
    container = false
    object    = false
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "%s"
  expiry = "%s"

  permissions {
    read    = true
    write   = true
    delete  = false
    list    = false
    add     = true
    create  = true
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, ipAddresses, startDate, endDate)
}

func (d StorageAccountSasDataSource) noPermissions(data acceptance.TestData, startDate string, endDate string, ipAddresses string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true
  ip_addresses      = "%s"
  signed_version    = "2019-10-10"

  resource_types {
    service   = true
    container = false
    object    = false
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "%s"
  expiry = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, ipAddresses, startDate, endDate)
}

func TestDataSourceStorageAccountSas_resourceTypesString(t *testing.T) {
	testCases := []struct {
		input    map[string]any
		expected string
	}{
		{map[string]any{"service": true}, "s"},
		{map[string]any{"container": true}, "c"},
		{map[string]any{"object": true}, "o"},
		{map[string]any{"service": true, "container": true, "object": true}, "sco"},
	}

	for _, test := range testCases {
		result := storage.BuildResourceTypesString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}

func TestDataSourceStorageAccountSas_servicesString(t *testing.T) {
	testCases := []struct {
		input    map[string]any
		expected string
	}{
		{map[string]any{"blob": true}, "b"},
		{map[string]any{"queue": true}, "q"},
		{map[string]any{"table": true}, "t"},
		{map[string]any{"file": true}, "f"},
		{map[string]any{"blob": true, "queue": true, "table": true, "file": true}, "bqtf"},
	}

	for _, test := range testCases {
		result := storage.BuildServicesString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}

func TestDataSourceStorageAccountSas_permissionsString(t *testing.T) {
	testCases := []struct {
		input    map[string]any
		expected string
	}{
		{map[string]any{"read": true}, "r"},
		{map[string]any{"write": true}, "w"},
		{map[string]any{"delete": true}, "d"},
		{map[string]any{"list": true}, "l"},
		{map[string]any{"add": true}, "a"},
		{map[string]any{"create": true}, "c"},
		{map[string]any{"update": true}, "u"},
		{map[string]any{"process": true}, "p"},
		{map[string]any{"tag": true}, "t"},
		{map[string]any{"filter": true}, "f"},
		{map[string]any{"read": true, "write": true, "add": true, "create": true}, "rwac"},
	}

	for _, test := range testCases {
		result := storage.BuildPermissionsString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}
