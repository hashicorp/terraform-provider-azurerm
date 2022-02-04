package recoveryservices_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RecoveryServicesVaultResource struct {
}

func TestAccRecoveryServicesVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

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

func TestAccRecoveryServicesVault_zrs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zrs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

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

func TestAccRecoveryServicesVault_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

func TestAccRecoveryServicesVault_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccRecoveryServicesVault_basicWithIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t RecoveryServicesVaultResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VaultID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.VaultsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func TestAccRecoveryServicesVault_encryptionWithKeyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkEncryptionWithKeyVaultKey(data, false, 0),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_turnOnEncryptionWithKeyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithIdentity(data),
		},
		data.ImportStep(),
		{
			Config: r.cmkEncryptionWithKeyVaultKey(data, false, 0),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_turnOffEncryptionWithKeyVaultKeyShouldHaveClearlyErrorMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkEncryptionWithKeyVaultKey(data, false, 0),
		},
		data.ImportStep(),
		{
			Config:      r.basicWithIdentity(data),
			ExpectError: regexp.MustCompile("once encryption with your own key has been enabled it's not possible to disable it"),
		},
	})
}

func TestAccRecoveryServicesVault_changeInfrastructureEncryptionEnabledShouldHaveClearlyErrorMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkEncryptionWithKeyVaultKey(data, false, 0),
		},
		data.ImportStep(),
		{
			Config:      r.cmkEncryptionWithKeyVaultKey(data, true, 0),
			ExpectError: regexp.MustCompile("once `infrastructure_encryption_enabled` has been set it's not possible to change it"),
		},
	})
}

func TestAccRecoveryServicesVault_encryptionWithInfrastructureEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkEncryptionWithKeyVaultKey(data, true, 0),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_switchEncryptionKeyVaultKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkEncryptionWithKeyVaultKey(data, false, 0),
		},
		data.ImportStep(),
		{
			Config: r.cmkEncryptionWithKeyVaultKey(data, false, 1),
		},
		data.ImportStep(),
	})
}

func (RecoveryServicesVaultResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) basicWithIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  identity {
    type = "SystemAssigned"
  }

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
  storage_mode_type   = "LocallyRedundant"
  tags = {
    ENV = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r RecoveryServicesVaultResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_recovery_services_vault" "import" {
  name                = azurerm_recovery_services_vault.test.name
  location            = azurerm_recovery_services_vault.test.location
  resource_group_name = azurerm_recovery_services_vault.test.resource_group_name
  sku                 = azurerm_recovery_services_vault.test.sku
}
`, r.basic(data))
}

func (r RecoveryServicesVaultResource) cmkEncryptionWithKeyVaultKey(data acceptance.TestData, enableInfraEncryptionState bool, keyIndex int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d"
  location = "%[2]s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  identity {
    type = "SystemAssigned"
  }

  soft_delete_enabled = true

  encryption {
    key_id                            = azurerm_key_vault_key.test[%[5]d].id
    use_system_assigned_identity      = true
    infrastructure_encryption_enabled = %[4]t
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                        = "acctest-key-vault-%[3]s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Decrypt",
      "Encrypt",
      "Delete",
      "Get",
      "List",
      "Purge",
      "UnwrapKey",
      "WrapKey",
      "Verify",
    ]
    secret_permissions = [
      "set",
    ]
  }
  lifecycle {
    ignore_changes = [soft_delete_enabled]
  }
}

resource "azurerm_key_vault_key" "test" {
  count        = 2
  name         = "acctest-key-vault-key-%[1]d${count.index}"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString, enableInfraEncryptionState, keyIndex)
}

func (RecoveryServicesVaultResource) zrs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false
  storage_mode_type   = "ZoneRedundant"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
