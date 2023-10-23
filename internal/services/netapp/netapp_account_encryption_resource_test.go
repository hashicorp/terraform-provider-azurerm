// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppAccountEncryptionResource struct{}

func TestAccNetAppAccountEncryption_cmkSystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkSystemAssigned(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption.0.key_vault_key_id").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppAccountEncryption_cmkUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkUserAssigned(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption.0.key_vault_key_id").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppAccountEncryption_updateKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account_encryption", "test")
	r := NetAppAccountEncryptionResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	regexInitialKey, err := regexp.Compile(`anfenckey\d+$`)
	if err != nil {
		t.Fatal(err)
	}

	regexNewKey, err := regexp.Compile(`.*anfenckey-new.*`)
	if err != nil {
		t.Fatal(err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.keyUpdate1(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption.0.key_vault_key_id").MatchesRegex(regexInitialKey),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyUpdate2(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption.0.key_vault_key_id").MatchesRegex(regexNewKey),
			),
		},
		data.ImportStep(),
	})
}

func (t NetAppAccountEncryptionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := netappaccounts.ParseNetAppAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.AccountClient.AccountsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Netapp Account (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r NetAppAccountEncryptionResource) cmkSystemAssigned(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
	name                            = "anfakv%[3]d"
	location                        = azurerm_resource_group.test.location
	resource_group_name             = azurerm_resource_group.test.name
	enabled_for_disk_encryption     = true
	enabled_for_deployment          = true
	enabled_for_template_deployment = true
	purge_protection_enabled        = true
	tenant_id                       = "%[4]s"

	sku_name = "standard"

	tags = {
		"CreatedOnDate" = "2022-07-08T23:50:21Z"
	}
}

resource "azurerm_key_vault_access_policy" "test-currentuser" {
	key_vault_id = azurerm_key_vault.test.id
	tenant_id = azurerm_netapp_account.test.identity.0.tenant_id
	object_id = data.azurerm_client_config.current.object_id

	key_permissions = [
		"Get",
		"Create",
		"Delete",
		"WrapKey",
		"UnwrapKey",
		"GetRotationPolicy",
		"SetRotationPolicy",
	]
}

resource "azurerm_key_vault_key" "test" {
	name         = "anfenckey%[3]d"
	key_vault_id = azurerm_key_vault.test.id
	key_type     = "RSA"
	key_size     = 2048

	key_opts = [
		"decrypt",
		"encrypt",
		"sign",
		"unwrapKey",
		"verify",
		"wrapKey",
	]

	depends_on = [
       azurerm_key_vault_access_policy.test-currentuser
	]
}

resource "azurerm_private_endpoint" "test" {
	name                   = "acctest-pe-akv-%[3]d"
	location               = azurerm_resource_group.test.location
	resource_group_name    = azurerm_resource_group.test.name
	subnet_id              = azurerm_subnet.test-non-delegated.id

	private_service_connection {
		name                          = "acctest-pe-sc-akv-%[3]d"
		private_connection_resource_id = azurerm_key_vault.test.id
		is_manual_connection          = false
		subresource_names             = ["Vault"]
	}

	tags = {
		CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
	}
}

resource "azurerm_netapp_account" "test" {
	name                = "acctest-NetAppAccount-%[3]d"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name

	identity {
		type = "SystemAssigned"
	}

	tags = {
		"CreatedOnDate" = "2022-07-08T23:50:21Z"
	}
}

resource "azurerm_key_vault_access_policy" "test-systemassigned" {
	key_vault_id = azurerm_key_vault.test.id
	tenant_id    = azurerm_netapp_account.test.identity.0.tenant_id
	object_id    = azurerm_netapp_account.test.identity.0.principal_id

	key_permissions = [
		"Get",
		"Encrypt",
		"Decrypt"
	]
}

resource "azurerm_netapp_account_encryption" "test" {
	netapp_account_id = azurerm_netapp_account.test.id

	system_assigned_identity_principal_id = azurerm_netapp_account.test.identity.0.principal_id

	encryption {
		key_vault_key_id = azurerm_key_vault_key.test.versionless_id
	}

	depends_on = [
		azurerm_key_vault_access_policy.test-systemassigned,
		azurerm_private_endpoint.test
	]
}
`, r.template(data), r.networkTemplate(data), data.RandomInteger, tenantID)
}

func (r NetAppAccountEncryptionResource) cmkUserAssigned(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_user_assigned_identity" "test" {
	name                = "user-assigned-identity-%[3]d"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name
  
	tags = {
	  CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
	}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
	name                            = "anfakv%[3]d"
	location                        = azurerm_resource_group.test.location
	resource_group_name             = azurerm_resource_group.test.name
	enabled_for_disk_encryption     = true
	enabled_for_deployment          = true
	enabled_for_template_deployment = true
	purge_protection_enabled        = true
	tenant_id                       = "%[4]s"
  
	sku_name = "standard"

	access_policy {
		tenant_id = "%[4]s"
		object_id = data.azurerm_client_config.current.object_id
	
		key_permissions = [
		  "Get",
		  "Create",
		  "Delete",
		  "WrapKey",
		  "UnwrapKey",
		  "GetRotationPolicy",
		  "SetRotationPolicy",
		]
	}

	access_policy {
		tenant_id = "%[4]s"
		object_id = azurerm_user_assigned_identity.test.principal_id
	
		key_permissions = [
		  "Get",
		  "Encrypt",
		  "Decrypt"
		]
	}

	tags = {
		"CreatedOnDate" = "2022-07-08T23:50:21Z"
	}
}

resource "azurerm_key_vault_key" "test" {
	name         = "anfenckey%[3]d"
	key_vault_id = azurerm_key_vault.test.id
	key_type     = "RSA"
	key_size     = 2048
  
	key_opts = [
	  "decrypt",
	  "encrypt",
	  "sign",
	  "unwrapKey",
	  "verify",
	  "wrapKey",
	]
}

resource "azurerm_private_endpoint" "test" {
	name                   = "acctest-pe-akv-%[3]d"
	location               = azurerm_resource_group.test.location
	resource_group_name    = azurerm_resource_group.test.name
	subnet_id              = azurerm_subnet.test-non-delegated.id

	private_service_connection {
	  name                          = "acctest-pe-sc-akv-%[3]d"
	  private_connection_resource_id = azurerm_key_vault.test.id
	  is_manual_connection          = false
	  subresource_names             = ["Vault"]
	}
	
	tags = {
		CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
	  }
}

resource "azurerm_netapp_account" "test" {
	name                = "acctest-NetAppAccount-%[3]d"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name
  
	identity {
		type = "UserAssigned"
		identity_ids = [
		  azurerm_user_assigned_identity.test.id
		]
	}

	tags = {
	  "CreatedOnDate" = "2022-07-08T23:50:21Z"
	}
}

resource "azurerm_netapp_account_encryption" "test" {
	netapp_account_id = azurerm_netapp_account.test.id
	
	user_assigned_identity_id = azurerm_user_assigned_identity.test.id

	encryption {
		key_vault_key_id = azurerm_key_vault_key.test.versionless_id
	}

	depends_on = [
		azurerm_private_endpoint.test
	]
}
`, r.template(data), r.networkTemplate(data), data.RandomInteger, tenantID)
}

func (r NetAppAccountEncryptionResource) keyUpdate1(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
	name                            = "anfakv%[3]d"
	location                        = azurerm_resource_group.test.location
	resource_group_name             = azurerm_resource_group.test.name
	enabled_for_disk_encryption     = true
	enabled_for_deployment          = true
	enabled_for_template_deployment = true
	purge_protection_enabled        = true
	tenant_id                       = "%[4]s"

	sku_name = "standard"

	tags = {
		"CreatedOnDate" = "2022-07-08T23:50:21Z"
	}
}

resource "azurerm_key_vault_access_policy" "test-currentuser" {
	key_vault_id = azurerm_key_vault.test.id
	tenant_id = azurerm_netapp_account.test.identity.0.tenant_id
	object_id = data.azurerm_client_config.current.object_id

	key_permissions = [
		"Get",
		"Create",
		"Delete",
		"WrapKey",
		"UnwrapKey",
		"GetRotationPolicy",
		"SetRotationPolicy",
	]
}

resource "azurerm_key_vault_key" "test" {
	name         = "anfenckey%[3]d"
	key_vault_id = azurerm_key_vault.test.id
	key_type     = "RSA"
	key_size     = 2048

	key_opts = [
		"decrypt",
		"encrypt",
		"sign",
		"unwrapKey",
		"verify",
		"wrapKey",
	]

	depends_on = [
       azurerm_key_vault_access_policy.test-currentuser
	]
}

resource "azurerm_key_vault_key" "test-new-key" {
	name         = "anfenckey-new%[3]d"
	key_vault_id = azurerm_key_vault.test.id
	key_type     = "RSA"
	key_size     = 2048

	key_opts = [
		"decrypt",
		"encrypt",
		"sign",
		"unwrapKey",
		"verify",
		"wrapKey",
	]

	depends_on = [
		azurerm_key_vault_key.test,
       azurerm_key_vault_access_policy.test-currentuser
	]
}

resource "azurerm_private_endpoint" "test" {
	name                   = "acctest-pe-akv-%[3]d"
	location               = azurerm_resource_group.test.location
	resource_group_name    = azurerm_resource_group.test.name
	subnet_id              = azurerm_subnet.test-non-delegated.id

	private_service_connection {
		name                          = "acctest-pe-sc-akv-%[3]d"
		private_connection_resource_id = azurerm_key_vault.test.id
		is_manual_connection          = false
		subresource_names             = ["Vault"]
	}

	tags = {
		CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
	}
}

resource "azurerm_netapp_account" "test" {
	name                = "acctest-NetAppAccount-%[3]d"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name

	identity {
		type = "SystemAssigned"
	}

	tags = {
		"CreatedOnDate" = "2022-07-08T23:50:21Z"
	}
}

resource "azurerm_key_vault_access_policy" "test-systemassigned" {
	key_vault_id = azurerm_key_vault.test.id
	tenant_id    = azurerm_netapp_account.test.identity.0.tenant_id
	object_id    = azurerm_netapp_account.test.identity.0.principal_id

	key_permissions = [
		"Get",
		"Encrypt",
		"Decrypt"
	]
}

resource "azurerm_netapp_account_encryption" "test" {
	netapp_account_id = azurerm_netapp_account.test.id

	system_assigned_identity_principal_id = azurerm_netapp_account.test.identity.0.principal_id

	encryption {
		key_vault_key_id = azurerm_key_vault_key.test.versionless_id
	}

	depends_on = [
		azurerm_key_vault_access_policy.test-systemassigned,
		azurerm_private_endpoint.test
	]
}

`, r.template(data), r.networkTemplate(data), data.RandomInteger, tenantID)
}

func (r NetAppAccountEncryptionResource) keyUpdate2(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "test" {
	name                            = "anfakv%[3]d"
	location                        = azurerm_resource_group.test.location
	resource_group_name             = azurerm_resource_group.test.name
	enabled_for_disk_encryption     = true
	enabled_for_deployment          = true
	enabled_for_template_deployment = true
	purge_protection_enabled        = true
	tenant_id                       = "%[4]s"

	sku_name = "standard"

	tags = {
		"CreatedOnDate" = "2022-07-08T23:50:21Z"
	}
}

resource "azurerm_key_vault_access_policy" "test-currentuser" {
	key_vault_id = azurerm_key_vault.test.id
	tenant_id = azurerm_netapp_account.test.identity.0.tenant_id
	object_id = data.azurerm_client_config.current.object_id

	key_permissions = [
		"Get",
		"Create",
		"Delete",
		"WrapKey",
		"UnwrapKey",
		"GetRotationPolicy",
		"SetRotationPolicy",
	]
}

resource "azurerm_key_vault_key" "test" {
	name         = "anfenckey%[3]d"
	key_vault_id = azurerm_key_vault.test.id
	key_type     = "RSA"
	key_size     = 2048

	key_opts = [
		"decrypt",
		"encrypt",
		"sign",
		"unwrapKey",
		"verify",
		"wrapKey",
	]

	depends_on = [
       azurerm_key_vault_access_policy.test-currentuser
	]
}

resource "azurerm_key_vault_key" "test-new-key" {
	name         = "anfenckey-new%[3]d"
	key_vault_id = azurerm_key_vault.test.id
	key_type     = "RSA"
	key_size     = 2048

	key_opts = [
		"decrypt",
		"encrypt",
		"sign",
		"unwrapKey",
		"verify",
		"wrapKey",
	]

	depends_on = [
		azurerm_key_vault_key.test,
       azurerm_key_vault_access_policy.test-currentuser
	]
}

resource "azurerm_private_endpoint" "test" {
	name                   = "acctest-pe-akv-%[3]d"
	location               = azurerm_resource_group.test.location
	resource_group_name    = azurerm_resource_group.test.name
	subnet_id              = azurerm_subnet.test-non-delegated.id

	private_service_connection {
		name                          = "acctest-pe-sc-akv-%[3]d"
		private_connection_resource_id = azurerm_key_vault.test.id
		is_manual_connection          = false
		subresource_names             = ["Vault"]
	}

	tags = {
		CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
	}
}

resource "azurerm_netapp_account" "test" {
	name                = "acctest-NetAppAccount-%[3]d"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name

	identity {
		type = "SystemAssigned"
	}

	tags = {
		"CreatedOnDate" = "2022-07-08T23:50:21Z"
	}
}

resource "azurerm_key_vault_access_policy" "test-systemassigned" {
	key_vault_id = azurerm_key_vault.test.id
	tenant_id    = azurerm_netapp_account.test.identity.0.tenant_id
	object_id    = azurerm_netapp_account.test.identity.0.principal_id

	key_permissions = [
		"Get",
		"Encrypt",
		"Decrypt"
	]
}

resource "azurerm_netapp_account_encryption" "test" {
	netapp_account_id = azurerm_netapp_account.test.id

	system_assigned_identity_principal_id = azurerm_netapp_account.test.identity.0.principal_id

	encryption {
		key_vault_key_id = azurerm_key_vault_key.test-new-key.versionless_id
	}

	depends_on = [
		azurerm_key_vault_access_policy.test-systemassigned,
		azurerm_private_endpoint.test
	]
}

`, r.template(data), r.networkTemplate(data), data.RandomInteger, tenantID)
}

func (NetAppAccountEncryptionResource) networkTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_virtual_network" "test" {
	name                = "acctest-VirtualNetwork-%[1]d"
	location            = azurerm_resource_group.test.location
	resource_group_name = azurerm_resource_group.test.name
	address_space       = ["10.6.0.0/16"]

	tags = {
		"CreatedOnDate"    = "2022-07-08T23:50:21Z",
		"SkipASMAzSecPack" = "true"
	}
}

resource "azurerm_subnet" "test-delegated" {
	name                 = "acctest-Delegated-Subnet-%[1]d"
	resource_group_name  = azurerm_resource_group.test.name
	virtual_network_name = azurerm_virtual_network.test.name
	address_prefixes     = ["10.6.1.0/24"]

	delegation {
		name = "testdelegation"

		service_delegation {
		name    = "Microsoft.Netapp/volumes"
		actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
		}
	}
}

resource "azurerm_subnet" "test-non-delegated" {
	name                 = "acctest-Non-Delegated-Subnet-%[1]d"
	resource_group_name  = azurerm_resource_group.test.name
	virtual_network_name = azurerm_virtual_network.test.name
	address_prefixes     = ["10.6.0.0/24"]
}
`, data.RandomInteger)
}

func (NetAppAccountEncryptionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }

    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipNRMSNSG"      = "true"
  }
}

`, data.RandomInteger, data.Locations.Primary)
}
