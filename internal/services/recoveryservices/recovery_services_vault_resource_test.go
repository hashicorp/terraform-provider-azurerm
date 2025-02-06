// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RecoveryServicesVaultResource struct{}

func TestAccRecoveryServicesVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cross_region_restore_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_ToggleCrossRegionRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithCrossRegionRestore(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cross_region_restore_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithCrossRegionRestore(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cross_region_restore_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithCrossRegionRestore(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cross_region_restore_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_basicWithCrossRegionRestoreAndWrongStorageType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basicWithCrossRegionRestoreAndWrongStorageType(data),
			ExpectError: regexp.MustCompile("cannot enable cross region restore when storage mode type is not GeoRedundant"),
		},
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

func TestAccRecoveryServicesVault_SystemAssignedIdentity(t *testing.T) {
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

func TestAccRecoveryServicesVault_Identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithSystemAssignedUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_immutabilityLocked(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// To test creation with `Locked`, it is irreversible.
			Config: r.basicWithImmutability(data, "Locked"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithImmutability(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithImmutability(data, "Locked"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_immutability(t *testing.T) {
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
			Config: r.basicWithImmutability(data, "Unlocked"),
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
		{
			Config: r.basicWithImmutability(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithImmutability(data, "Unlocked"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithImmutability(data, "Locked"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_softDelete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.softDeleteDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.softDeleteDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.softDeleteDisabledWithEncryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_storageModeType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageModeTypeDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.storageModeTypeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_crossRegionRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossRegionRestoreDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.crossRegionRestoreEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_CrossRegionRestoreWithEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossRegionRestoreEnabledWithEncryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_sku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.skuStandard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.skuRS0(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// This step is to ensure `sku` and `encryption` can be updated together when `encryption` was not specified in the existing resource
			Config: r.skuWithEncryption(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_updateSkuWithExistingEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.skuWithEncryption(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.skuWithEncryption(data, "RS0"),
			ExpectError: regexp.MustCompile("`sku` cannot be changed when encryption with your own key has been enabled"),
		},
	})
}

func (t RecoveryServicesVaultResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := vaults.ParseVaultID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RecoveryServices.VaultsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Recovery Service (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
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

func TestAccRecoveryServicesVault_encryptionWithUserAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cmkEncryptionWithUserAssignedIdentity(data, false, 0),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_TogglePublicNetworkAccessEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithPublicNetworkAccessEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithPublicNetworkAccessEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithPublicNetworkAccessEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_basicWithClassicVmwareReplicateEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithClassicVmwareReplicateEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("classic_vmware_replication_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRecoveryServicesVault_basicWithMonitorDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_recovery_services_vault", "test")
	r := RecoveryServicesVaultResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithMonitor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
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

func (RecoveryServicesVaultResource) basicWithPublicNetworkAccessEnabled(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                          = "acctest-Vault-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  sku                           = "Standard"
  public_network_access_enabled = %t

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enabled)
}

func (RecoveryServicesVaultResource) basicWithCrossRegionRestore(data acceptance.TestData, enable bool) string {
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

  cross_region_restore_enabled = %t

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enable)
}

func (RecoveryServicesVaultResource) basicWithCrossRegionRestoreAndWrongStorageType(data acceptance.TestData) string {
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

  storage_mode_type = "ZoneRedundant"

  cross_region_restore_enabled = true

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

func (RecoveryServicesVaultResource) basicWithUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary)
}

func (RecoveryServicesVaultResource) basicWithSystemAssignedUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
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
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary)
}

func (RecoveryServicesVaultResource) basicWithImmutability(data acceptance.TestData, immutability string) string {
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
  immutability        = "%s"

  soft_delete_enabled = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, immutability)
}

func (RecoveryServicesVaultResource) basicWithMonitor(data acceptance.TestData) string {
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

  monitoring {
    alerts_for_all_job_failures_enabled            = false
    alerts_for_critical_operation_failures_enabled = false
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
      purge_soft_delete_on_destroy       = true
      purge_soft_deleted_keys_on_destroy = false
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
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
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

func (r RecoveryServicesVaultResource) cmkEncryptionWithUserAssignedIdentity(data acceptance.TestData, enableInfraEncryptionState bool, keyIndex int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = true
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-identity-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  encryption {
    key_id                            = azurerm_key_vault_key.test[%[5]d].id
    use_system_assigned_identity      = false
    user_assigned_identity_id         = azurerm_user_assigned_identity.test.id
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
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = [
      "Get",
      "List",
      "WrapKey",
      "UnwrapKey",
      "GetRotationPolicy"
    ]
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

func (RecoveryServicesVaultResource) softDeleteDefault(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) softDeleteDisabledWithEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = true
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d"
  location = "%[2]s"
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
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctest-key-vault-key"
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

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  soft_delete_enabled = false

  encryption {
    key_id                            = azurerm_key_vault_key.test.id
    infrastructure_encryption_enabled = false
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (RecoveryServicesVaultResource) softDeleteDisabled(data acceptance.TestData) string {
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

func (RecoveryServicesVaultResource) crossRegionRestoreDefault(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) crossRegionRestoreEnabled(data acceptance.TestData) string {
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

  cross_region_restore_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) crossRegionRestoreEnabledWithEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = true
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d"
  location = "%[2]s"
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
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctest-key-vault-key"
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

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  cross_region_restore_enabled = true

  encryption {
    key_id                            = azurerm_key_vault_key.test.id
    infrastructure_encryption_enabled = false
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (RecoveryServicesVaultResource) storageModeTypeDefault(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) storageModeTypeUpdated(data acceptance.TestData) string {
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

  storage_mode_type = "ZoneRedundant"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) skuStandard(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) skuRS0(data acceptance.TestData) string {
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
  sku                 = "RS0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RecoveryServicesVaultResource) skuWithEncryption(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = true
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%[1]d"
  location = "%[2]s"
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
      "GetRotationPolicy"
    ]
    secret_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctest-key-vault-key"
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

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-Vault-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "%[4]s"

  encryption {
    key_id                            = azurerm_key_vault_key.test.id
    infrastructure_encryption_enabled = false
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, sku)
}

func (RecoveryServicesVaultResource) basicWithClassicVmwareReplicateEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-recovery-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                               = "acctest-Vault-%d"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  sku                                = "Standard"
  classic_vmware_replication_enabled = true
  soft_delete_enabled                = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
