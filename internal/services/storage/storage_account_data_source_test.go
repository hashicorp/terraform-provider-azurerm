// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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
				check.That(data.ResourceName).Key("public_network_access").HasValue("Enabled"),
				check.That(data.ResourceName).Key("allowed_copy_scope").HasValue(""),
				check.That(data.ResourceName).Key("cross_tenant_replication_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("default_to_oauth_authentication").HasValue("false"),
				check.That(data.ResourceName).Key("edge_zone").HasValue(""),
				check.That(data.ResourceName).Key("large_file_share_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("local_user_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("provisioned_billing_model_version").HasValue(""),
				check.That(data.ResourceName).Key("sftp_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("shared_access_key_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("customer_managed_key.#").HasValue("0"),
				check.That(data.ResourceName).Key("immutability_policy.#").HasValue("0"),
				check.That(data.ResourceName).Key("network_rules.#").HasValue("0"),
				check.That(data.ResourceName).Key("routing.#").HasValue("0"),
				check.That(data.ResourceName).Key("sas_policy.#").HasValue("0"),
				check.That(data.ResourceName).Key("blob_properties.#").HasValue("1"),
				check.That(data.ResourceName).Key("blob_properties.0.versioning_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("blob_properties.0.change_feed_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("share_properties.#").HasValue("1"),
				check.That(data.ResourceName).Key("share_properties.0.retention_policy.#").HasValue("0"),
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

func TestAccDataSourceStorageAccount_withEncryptionKey_Service(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.encryptionKeyWithDataSource(data, "Service"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("table_encryption_key_type").HasValue("Service"),
				check.That(data.ResourceName).Key("queue_encryption_key_type").HasValue("Service"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_withEncryptionKey_Account(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.encryptionKeyWithDataSource(data, "Account"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("table_encryption_key_type").HasValue("Account"),
				check.That(data.ResourceName).Key("queue_encryption_key_type").HasValue("Account"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_withInfrastructureEncryptionEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.infrastructureEncryptionWithDataSource(data, "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("infrastructure_encryption_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_withInfrastructureEncryptionDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
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

func TestAccDataSourceStorageAccount_azureFilesAuthentication(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.azureFilesAuthenticationAD(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("azure_files_authentication.0.directory_type").HasValue("AD"),
				check.That(data.ResourceName).Key("azure_files_authentication.0.active_directory.0.storage_sid").HasValue("S-1-5-21-2400535526-2334094090-2402026252-0012"),
				check.That(data.ResourceName).Key("azure_files_authentication.0.active_directory.0.domain_name").HasValue("adtest.com"),
				check.That(data.ResourceName).Key("azure_files_authentication.0.active_directory.0.domain_sid").HasValue("S-1-5-21-2400535526-2334094090-2402026252-0012"),
				check.That(data.ResourceName).Key("azure_files_authentication.0.active_directory.0.domain_guid").HasValue("aebfc118-9fa9-4732-a21f-d98e41a77ae1"),
				check.That(data.ResourceName).Key("azure_files_authentication.0.active_directory.0.forest_name").HasValue("adtest.com"),
				check.That(data.ResourceName).Key("azure_files_authentication.0.active_directory.0.netbios_domain_name").HasValue("adtest.com"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("allowed_copy_scope").HasValue("AAD"),
				check.That(data.ResourceName).Key("cross_tenant_replication_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("default_to_oauth_authentication").HasValue("true"),
				check.That(data.ResourceName).Key("large_file_share_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("local_user_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("shared_access_key_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("network_rules.0.default_action").HasValue("Deny"),
				check.That(data.ResourceName).Key("network_rules.0.ip_rules.#").HasValue("1"),
				check.That(data.ResourceName).Key("routing.0.choice").HasValue("InternetRouting"),
				check.That(data.ResourceName).Key("routing.0.publish_internet_endpoints").HasValue("true"),
				check.That(data.ResourceName).Key("routing.0.publish_microsoft_endpoints").HasValue("false"),
				check.That(data.ResourceName).Key("sas_policy.0.expiration_action").HasValue("Log"),
				check.That(data.ResourceName).Key("sas_policy.0.expiration_period").HasValue("1.15:5:05"),
				check.That(data.ResourceName).Key("blob_properties.0.versioning_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("blob_properties.0.change_feed_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("blob_properties.0.delete_retention_policy.0.days").HasValue("7"),
				check.That(data.ResourceName).Key("blob_properties.0.cors_rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("blob_properties.0.cors_rule.0.maximum_age_in_seconds").HasValue("500"),
				check.That(data.ResourceName).Key("share_properties.0.retention_policy.0.days").HasValue("90"),
				check.That(data.ResourceName).Key("share_properties.0.cors_rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("share_properties.0.cors_rule.0.maximum_age_in_seconds").HasValue("1000"),
				check.That(data.ResourceName).Key("share_properties.0.smb.0.versions.#").HasValue("1"),
				check.That(data.ResourceName).Key("immutability_policy.0.state").HasValue("Unlocked"),
				check.That(data.ResourceName).Key("immutability_policy.0.period_since_creation_in_days").HasValue("1"),
				check.That(data.ResourceName).Key("immutability_policy.0.allow_protected_append_writes").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_customerManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.customerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("customer_managed_key.#").HasValue("1"),
				check.That(data.ResourceName).Key("customer_managed_key.0.key_vault_key_id").MatchesOtherKey(
					check.That("azurerm_key_vault_key.test").Key("id"),
				),
				check.That(data.ResourceName).Key("customer_managed_key.0.user_assigned_identity_id").MatchesOtherKey(
					check.That("azurerm_user_assigned_identity.test").Key("id"),
				),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_edgeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.edgeZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("edge_zone").MatchesOtherKey(
					check.That("azurerm_storage_account.test").Key("edge_zone"),
				),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_provisionedBillingModelVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.provisionedBillingModelVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("provisioned_billing_model_version").HasValue("V2"),
			),
		},
	})
}

func TestAccDataSourceStorageAccount_sftp(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: StorageAccountDataSource{}.sftp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sftp_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("is_hns_enabled").HasValue("true"),
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

func (d StorageAccountDataSource) azureFilesAuthenticationAD(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  azure_files_authentication {
    directory_type = "AD"
    active_directory {
      storage_sid         = "S-1-5-21-2400535526-2334094090-2402026252-0012"
      domain_name         = "adtest.com"
      domain_sid          = "S-1-5-21-2400535526-2334094090-2402026252-0012"
      domain_guid         = "aebfc118-9fa9-4732-a21f-d98e41a77ae1"
      forest_name         = "adtest.com"
      netbios_domain_name = "adtest.com"
    }
  }
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (d StorageAccountDataSource) complete(data acceptance.TestData) string {
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

  location                         = azurerm_resource_group.test.location
  account_kind                     = "StorageV2"
  account_tier                     = "Standard"
  account_replication_type         = "LRS"
  is_hns_enabled                   = false
  local_user_enabled               = false
  cross_tenant_replication_enabled = true
  default_to_oauth_authentication  = true
  shared_access_key_enabled        = false
  large_file_share_enabled         = true
  allowed_copy_scope               = "AAD"

  network_rules {
    default_action = "Deny"
    bypass         = ["Metrics"]
    ip_rules       = ["127.0.0.1"]
  }

  routing {
    choice                      = "InternetRouting"
    publish_internet_endpoints  = true
    publish_microsoft_endpoints = false
  }

  sas_policy {
    expiration_action = "Log"
    expiration_period = "1.15:5:05"
  }

  blob_properties {
    cors_rule {
      allowed_origins    = ["http://www.example.com"]
      allowed_methods    = ["GET"]
      allowed_headers    = ["x-tempo-*"]
      exposed_headers    = ["x-tempo-*"]
      max_age_in_seconds = 500
    }

    delete_retention_policy {
      days = 7
    }

    versioning_enabled       = true
    change_feed_enabled      = true
    last_access_time_enabled = true
  }

  share_properties {
    cors_rule {
      allowed_origins    = ["http://www.test.com"]
      allowed_methods    = ["PUT"]
      allowed_headers    = ["*"]
      exposed_headers    = ["x-tempo-*"]
      max_age_in_seconds = 1000
    }

    retention_policy {
      days = 90
    }

    smb {
      versions                        = ["SMB3.0"]
      authentication_types            = ["NTLMv2"]
      kerberos_ticket_encryption_type = ["AES-256"]
    }
  }

  immutability_policy {
    period_since_creation_in_days = 1
    state                         = "Unlocked"
    allow_protected_append_writes = false
  }
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (d StorageAccountDataSource) edgeZone(data acceptance.TestData) string {
	// @tombuildsstuff: WestUS has an edge zone available - so hard-code to that for now
	data.Locations.Primary = "westus"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Premium"
  account_replication_type = "LRS"
  edge_zone                = data.azurerm_extended_locations.test.extended_locations[0]

  tags = {
    environment = "production"
  }
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (d StorageAccountDataSource) provisionedBillingModelVersion(data acceptance.TestData) string {
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
  provisioned_billing_model_version = "V2"
  account_replication_type          = "LRS"
  account_kind                      = "FileStorage"
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (d StorageAccountDataSource) customerManagedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  customer_managed_key {
    key_vault_key_id          = azurerm_key_vault_key.test.id
    user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }

  infrastructure_encryption_enabled = true
  table_encryption_key_type         = "Account"
  queue_encryption_key_type         = "Account"
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, StorageAccountResource{}.cmkTemplate(data), data.RandomString)
}

func (d StorageAccountDataSource) sftp(data acceptance.TestData) string {
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
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
  sftp_enabled             = true
}

data "azurerm_storage_account" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
