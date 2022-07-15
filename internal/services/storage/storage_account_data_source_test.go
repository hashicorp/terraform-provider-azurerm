package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StorageAccountDataSource struct{}

func TestAccDataSourceStorageAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.basic(data),
		},
		{
			Config: StorageAccountDataSource{}.basicWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("account_tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("account_replication_type").HasValue("LRS"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("production"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_withWriteLock(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.basicWriteLock(data),
		},
		{
			Config: StorageAccountDataSource{}.basicWriteLockWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("account_tier").HasValue("Standard"),
				check.That(data.ResourceName).Key("account_replication_type").HasValue("LRS"),
				check.That(data.ResourceName).Key("primary_connection_string").IsEmpty(),
				check.That(data.ResourceName).Key("secondary_connection_string").IsEmpty(),
				check.That(data.ResourceName).Key("primary_blob_connection_string").IsEmpty(),
				check.That(data.ResourceName).Key("secondary_blob_connection_string").IsEmpty(),
				check.That(data.ResourceName).Key("primary_access_key").IsEmpty(),
				check.That(data.ResourceName).Key("secondary_access_key").IsEmpty(),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_withEncryptionKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.encryptionKeyWithDataSource(data, "Service"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("table_encryption_key_type").HasValue("Service"),
				check.That(data.ResourceName).Key("queue_encryption_key_type").HasValue("Service"),
			),
		},
		{
			Config: StorageAccountDataSource{}.encryptionKeyWithDataSource(data, "Account"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("table_encryption_key_type").HasValue("Account"),
				check.That(data.ResourceName).Key("queue_encryption_key_type").HasValue("Account"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_withInfrastructureEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.infrastructureEncryptionWithDataSource(data, "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("infrastructure_encryption_enabled").HasValue("true"),
			),
		},
		{
			Config: StorageAccountDataSource{}.infrastructureEncryptionWithDataSource(data, "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("infrastructure_encryption_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_systemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.0").IsSet(),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_systemAssignedUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.systemAssignedUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.0").IsSet(),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
	})
}

func (d StorageAccountDataSource) basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (d StorageAccountDataSource) basicWriteLock(data acceptance.TestData) string {
	template := d.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_storage_account.test.id
  lock_level = "ReadOnly"
}
`, template, data.RandomInteger)
}

func (d StorageAccountDataSource) basicWithDataSource(data acceptance.TestData) string {
	config := d.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, config)
}

func (d StorageAccountDataSource) basicWriteLockWithDataSource(data acceptance.TestData) string {
	config := d.basicWriteLock(data)
	return fmt.Sprintf(`
%s

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, config)
}

func (d StorageAccountDataSource) encryptionKeyWithDataSource(data acceptance.TestData, t string) string {
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

  location                  = azurerm_resource_group.test.location
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  table_encryption_key_type = %q
  queue_encryption_key_type = %q
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, t, t)
}

func (d StorageAccountDataSource) infrastructureEncryptionWithDataSource(data acceptance.TestData, t string) string {
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

  location                          = azurerm_resource_group.test.location
  account_tier                      = "Standard"
  account_replication_type          = "LRS"
  infrastructure_encryption_enabled = %s
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, t)
}

func (d StorageAccountDataSource) identityTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (d StorageAccountDataSource) systemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, d.identityTemplate(data), data.RandomString)
}

func (d StorageAccountDataSource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, d.identityTemplate(data), data.RandomString)
}

func (d StorageAccountDataSource) systemAssignedUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, d.identityTemplate(data), data.RandomString)
}
