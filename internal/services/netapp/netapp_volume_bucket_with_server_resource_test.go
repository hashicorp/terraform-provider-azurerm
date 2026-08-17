// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

type NetAppVolumeBucketWithServerResource struct{}

func TestAccNetAppVolumeBucketWithServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket_with_server", "test")
	r := NetAppVolumeBucketWithServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("server.0.certificate_pem"),
	})
}

func TestAccNetAppVolumeBucketWithServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket_with_server", "test")
	r := NetAppVolumeBucketWithServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("server.0.certificate_pem"),
	})
}

func TestAccNetAppVolumeBucketWithServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket_with_server", "test")
	r := NetAppVolumeBucketWithServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("server.0.certificate_pem"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("server.0.certificate_pem"),
	})
}

func TestAccNetAppVolumeBucketWithServer_withKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket_with_server", "test")
	r := NetAppVolumeBucketWithServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withKeyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeBucketWithServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket_with_server", "test")
	r := NetAppVolumeBucketWithServerResource{}

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

func (t NetAppVolumeBucketWithServerResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := buckets.ParseBucketID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.BucketsClient.Get(ctx, *id)
	if err != nil {
		if resp.HttpResponse != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (NetAppVolumeBucketWithServerResource) basic(data acceptance.TestData) string {
	template := NetAppVolumeBucketWithServerResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket_with_server" "test" {
  name      = "acctest-bucket-%[2]d"
  volume_id = azurerm_netapp_volume.test.id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  server {
    fqdn            = local.bucket_fqdn
    certificate_pem = base64encode("${tls_self_signed_cert.test.cert_pem}${tls_private_key.test.private_key_pem}")
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeBucketWithServerResource) complete(data acceptance.TestData) string {
	template := NetAppVolumeBucketWithServerResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket_with_server" "test" {
  name        = "acctest-bucket-%[2]d"
  volume_id   = azurerm_netapp_volume.test.id
  permissions = "ReadWrite"

  file_system_nfs_user {
    group_id = 2000
    user_id  = 2000
  }

  server {
    fqdn            = local.bucket_fqdn
    certificate_pem = base64encode("${tls_self_signed_cert.test.cert_pem}${tls_private_key.test.private_key_pem}")
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeBucketWithServerResource) withKeyVault(data acceptance.TestData) string {
	template := NetAppVolumeBucketWithServerResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "cert" {
  name                       = "kvcert%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  rbac_authorization_enabled = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = false
}

resource "azurerm_key_vault" "creds" {
  name                       = "kvcred%[3]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  rbac_authorization_enabled = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = false
}

resource "azurerm_key_vault_access_policy" "deployer_cert" {
  key_vault_id = azurerm_key_vault.cert.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  certificate_permissions = ["Get", "List", "Create", "Import", "Update", "Delete", "Purge", "Recover"]
  secret_permissions      = ["Get", "List", "Set", "Delete", "Purge", "Recover"]
}

resource "azurerm_key_vault_access_policy" "anf_cert" {
  key_vault_id = azurerm_key_vault.cert.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_netapp_account.test.identity[0].principal_id

  certificate_permissions = ["Get", "List", "Update", "Create", "Import", "ManageContacts", "GetIssuers", "ListIssuers", "SetIssuers", "DeleteIssuers"]
  secret_permissions      = ["Get", "List", "Set", "Delete"]
}

resource "azurerm_key_vault_access_policy" "anf_creds" {
  key_vault_id = azurerm_key_vault.creds.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_netapp_account.test.identity[0].principal_id

  secret_permissions = ["Get", "List", "Set", "Delete"]
}

resource "azurerm_key_vault_certificate" "bucket" {
  name         = "acctbucketcert%[2]d"
  key_vault_id = azurerm_key_vault.cert.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage          = ["digitalSignature", "keyEncipherment"]
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]
      subject            = "CN=${local.bucket_fqdn}"

      subject_alternative_names {
        dns_names = [local.bucket_fqdn]
      }

      validity_in_months = 12
    }
  }

  depends_on = [azurerm_key_vault_access_policy.deployer_cert]
}

resource "azurerm_netapp_volume_bucket_with_server" "test" {
  name      = "acctest-bucket-kv-%[2]d"
  volume_id = azurerm_netapp_volume.test.id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  server {
    fqdn = local.bucket_fqdn
  }

  key_vault {
    certificate_key_vault_uri = azurerm_key_vault.cert.vault_uri
    certificate_name          = azurerm_key_vault_certificate.bucket.name
    credentials_key_vault_uri = azurerm_key_vault.creds.vault_uri
    credentials_secret_name   = "acctbucketcreds%[2]d"
  }

  depends_on = [
    azurerm_key_vault_access_policy.anf_cert,
    azurerm_key_vault_access_policy.anf_creds,
  ]
}
`, template, data.RandomInteger, data.RandomString)
}

func (r NetAppVolumeBucketWithServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket_with_server" "import" {
  name      = azurerm_netapp_volume_bucket_with_server.test.name
  volume_id = azurerm_netapp_volume_bucket_with_server.test.volume_id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  server {
    fqdn            = local.bucket_fqdn
    certificate_pem = base64encode("${tls_self_signed_cert.test.cert_pem}${tls_private_key.test.private_key_pem}")
  }
}
`, r.basic(data))
}

func (NetAppVolumeBucketWithServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    netapp {
      prevent_volume_destruction             = false
      delete_backups_on_backup_vault_destroy = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

locals {
  bucket_fqdn = "acctbucket-%[1]d.example.internal"
}

resource "tls_private_key" "test" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test" {
  private_key_pem = tls_private_key.test.private_key_pem

  subject {
    common_name = local.bucket_fqdn
  }

  dns_names = [local.bucket_fqdn]

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "CreatedOnDate"    = "2026-01-15T00-00-00Z",
    "SkipASMAzSecPack" = "true",
    "SkipNRMSNSG"      = "true"
  }
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.99.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-DelegatedSubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.99.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

# System-assigned identity is required so that Azure NetApp Files can read
# the bucket server certificate from Key Vault and write the generated bucket
# credentials to Key Vault when the key_vault block is used on the bucket.
resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Auto"
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%[1]d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
  protocols           = ["NFSv3"]
}
`, data.RandomInteger, data.Locations.Primary)
}
