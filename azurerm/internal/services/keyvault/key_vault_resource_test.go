package keyvault_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultResource struct {
}

func TestAccKeyVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("standard"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault"),
		},
	})
}

func TestAccKeyVault_networkAcls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkAcls(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkAclsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_networkAclsAllowed(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.networkAclsAllowed(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_accessPolicyUpperLimit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.accessPolicyUpperLimit(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckKeyVaultDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccKeyVault_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckKeyVaultDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccKeyVault_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("access_policy.0.application_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("access_policy.0.key_permissions.0").HasValue("create"),
				check.That(data.ResourceName).Key("access_policy.0.secret_permissions.0").HasValue("set"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("access_policy.0.key_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("access_policy.0.secret_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("enabled_for_deployment").HasValue("true"),
				check.That(data.ResourceName).Key("enabled_for_disk_encryption").HasValue("true"),
				check.That(data.ResourceName).Key("enabled_for_template_deployment").HasValue("true"),
				check.That(data.ResourceName).Key("enable_rbac_authorization").HasValue("true"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Staging"),
			),
		},
		{
			Config: r.noAccessPolicyBlocks(data),
			Check: resource.ComposeTestCheckFunc(
				// There are no access_policy blocks in this configuration
				// at all, which means to ignore any existing policies and
				// so the one created in previous steps is still present.
				check.That(data.ResourceName).Key("access_policy.#").HasValue("1"),
			),
		},
		{
			Config: r.accessPolicyExplicitZero(data),
			Check: resource.ComposeTestCheckFunc(
				// This config explicitly sets access_policy = [], which
				// means to delete any existing policies.
				check.That(data.ResourceName).Key("access_policy.#").HasValue("0"),
			),
		},
	})
}

func TestAccKeyVault_upgradeSKU(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("standard"),
			),
		},
		data.ImportStep(), {
			Config: r.basicPremiumSKU(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("premium"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_updateContacts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateContacts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_justCert(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.justCert(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("access_policy.0.certificate_permissions.0").HasValue("get"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_softDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.softDelete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_softDeleteRecovery(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// create it regularly
			Config: r.softDelete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			// delete the key vault
			Config: r.softDeleteAbsent(data),
		},
		{
			// attempting to re-create it requires recovery, which is enabled by default
			Config: r.softDelete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_softDeleteRecoveryDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// create it regularly
			Config: r.softDeleteRecoveryDisabled(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			// delete the key vault
			Config: r.softDeleteAbsent(data),
		},
		{
			// attempting to re-create it requires recovery, which is enabled by default
			Config:      r.softDeleteRecoveryDisabled(data),
			ExpectError: regexp.MustCompile("An existing soft-deleted Key Vault exists with the Name"),
		},
	})
}

func TestAccKeyVault_purgeProtectionEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.purgeProtection(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_purgeProtectionAndSoftDeleteEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.purgeProtectionAndSoftDelete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_purgeProtectionViaUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.purgeProtection(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.purgeProtection(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVault_purgeProtectionAttemptToDisable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.purgeProtection(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("purge_protection_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config:      r.purgeProtection(data, false),
			ExpectError: regexp.MustCompile("once Purge Protection has been Enabled it's not possible to disable it"),
		},
	})
}

func TestAccKeyVault_deletePolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.noPolicy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("access_policy.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (t KeyVaultResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["vaults"]

	resp, err := clients.KeyVault.VaultsClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading Key Vault (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckKeyVaultDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		vaultName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for vault: %s", vaultName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Delete(ctx, resourceGroup, vaultName)
		if err != nil {
			if response.WasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Delete on keyVaultClient: %+v", err)
		}

		return nil
	}
}

func (KeyVaultResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "managecontacts",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r KeyVaultResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "import" {
  name                       = azurerm_key_vault.test.name
  location                   = azurerm_key_vault.test.location
  resource_group_name        = azurerm_key_vault.test.resource_group_name
  tenant_id                  = azurerm_key_vault.test.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}
`, r.basic(data))
}

func (KeyVaultResource) networkAclsTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test_a" {
  name                 = "acctestsubneta%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.KeyVault"]
}

resource "azurerm_subnet" "test_b" {
  name                 = "acctestsubnetb%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.4.0/24"
  service_endpoints    = ["Microsoft.KeyVault"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KeyVaultResource) networkAcls(data acceptance.TestData) string {
	template := KeyVaultResource{}.networkAclsTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }

  network_acls {
    default_action             = "Deny"
    bypass                     = "None"
    virtual_network_subnet_ids = [azurerm_subnet.test_a.id, azurerm_subnet.test_b.id]
  }
}
`, template, data.RandomInteger)
}

func (r KeyVaultResource) networkAclsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }

  network_acls {
    default_action             = "Allow"
    bypass                     = "AzureServices"
    ip_rules                   = ["123.0.0.102/32"]
    virtual_network_subnet_ids = [azurerm_subnet.test_a.id]
  }
}
`, r.networkAclsTemplate(data), data.RandomInteger)
}

func (r KeyVaultResource) networkAclsAllowed(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }

  network_acls {
    default_action = "Allow"
    bypass         = "AzureServices"
  }
}
`, r.networkAclsTemplate(data), data.RandomInteger)
}

func (KeyVaultResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
    ]
  }

  enabled_for_deployment          = true
  enabled_for_disk_encryption     = true
  enabled_for_template_deployment = true
  enable_rbac_authorization       = true

  tags = {
    environment = "Staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) noAccessPolicyBlocks(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  enabled_for_deployment          = true
  enabled_for_disk_encryption     = true
  enabled_for_template_deployment = true
  enable_rbac_authorization       = true

  tags = {
    environment = "Staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) accessPolicyExplicitZero(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy = []

  enabled_for_deployment          = true
  enabled_for_disk_encryption     = true
  enabled_for_template_deployment = true
  enable_rbac_authorization       = true

  tags = {
    environment = "Staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id      = data.azurerm_client_config.current.tenant_id
    object_id      = data.azurerm_client_config.current.object_id
    application_id = data.azurerm_client_config.current.client_id

    certificate_permissions = [
      "get",
    ]

    key_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
    ]
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) justCert(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "get",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) accessPolicyUpperLimit(data acceptance.TestData) string {
	var storageAccountConfigs string
	var accessPoliciesConfigs string

	for i := 1; i <= 20; i++ {
		storageAccountConfigs += testAccKeyVault_generateStorageAccountConfigs(i, data.RandomString)
		accessPoliciesConfigs += testAccKeyVault_generateAccessPolicyConfigs(i)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  %s
}

%s
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, accessPoliciesConfigs, storageAccountConfigs)
}

func testAccKeyVault_generateStorageAccountConfigs(accountNum int, rs string) string {
	return fmt.Sprintf(`
resource "azurerm_storage_account" "test%d" {
  name                     = "testsa%s%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "testing"
  }
}
`, accountNum, rs, accountNum)
}

func testAccKeyVault_generateAccessPolicyConfigs(accountNum int) string {
	// due to a weird terraform fmt issue where:
	//   "${azurerm_storage_account.test%d.identity.0.principal_id}"
	// becomes
	//   "${azurerm_storage_account.test % d.identity.0.principal_id}"
	//
	// lets inject this separately so we can run terrafmt on this file

	oid := fmt.Sprintf("${azurerm_storage_account.test%d.identity.0.principal_id}", accountNum)

	return fmt.Sprintf(`
access_policy {
  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = "%s"

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}
`, oid)
}

func (KeyVaultResource) purgeProtection(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                     = "vault%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = "%t"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enabled)
}

func (KeyVaultResource) softDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) softDeleteAbsent(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults = false
    }
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (KeyVaultResource) softDeleteRecoveryDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) purgeProtectionAndSoftDelete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) noPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) updateContacts(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-kv-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "managecontacts",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }

  contact {
    email = "example@example.com"
    name  = "example"
    phone = "01234567890"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KeyVaultResource) basicPremiumSKU(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "vault%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "managecontacts",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
