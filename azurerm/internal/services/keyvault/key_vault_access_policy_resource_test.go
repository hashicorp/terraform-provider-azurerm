package keyvault_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultAccessPolicyResource struct {
}

func TestAccKeyVaultAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.0").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.1").HasValue("Set"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultAccessPolicy_mixedCasePermissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicMixedCase(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.0").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.1").HasValue("set"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.0").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.1").HasValue("Set"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_access_policy"),
		},
	})
}

func TestAccKeyVaultAccessPolicy_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test_with_application_id")
	r := KeyVaultAccessPolicyResource{}
	resourceName2 := "azurerm_key_vault_access_policy.test_no_application_id"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("Create"),
				check.That(data.ResourceName).Key("key_permissions.1").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.0").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.1").HasValue("Delete"),
				check.That(data.ResourceName).Key("certificate_permissions.0").HasValue("Create"),
				check.That(data.ResourceName).Key("certificate_permissions.1").HasValue("Delete"),
				acceptance.TestCheckResourceAttr(resourceName2, "key_permissions.0", "List"),
				acceptance.TestCheckResourceAttr(resourceName2, "key_permissions.1", "Encrypt"),
				acceptance.TestCheckResourceAttr(resourceName2, "secret_permissions.0", "List"),
				acceptance.TestCheckResourceAttr(resourceName2, "secret_permissions.1", "Delete"),
				acceptance.TestCheckResourceAttr(resourceName2, "certificate_permissions.0", "List"),
				acceptance.TestCheckResourceAttr(resourceName2, "certificate_permissions.1", "Delete"),
			),
		},
		data.ImportStep(),
		{
			ResourceName:      resourceName2,
			ImportState:       true,
			ImportStateVerify: true,
		},
	})
}

func TestAccKeyVaultAccessPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.0").HasValue("Get"),
				check.That(data.ResourceName).Key("secret_permissions.1").HasValue("Set"),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("List"),
				check.That(data.ResourceName).Key("key_permissions.1").HasValue("Encrypt"),
			),
		},
	})
}

func TestAccKeyVaultAccessPolicy_nonExistentVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:             r.nonExistentVault(data),
			ExpectNonEmptyPlan: true,
			ExpectError:        regexp.MustCompile(`retrieving Key Vault`),
		},
	})
}

func (t KeyVaultAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	objectId := id.Path["objectId"]
	applicationId := id.Path["applicationId"]

	resp, err := clients.KeyVault.VaultsClient.Get(ctx, resGroup, vaultName)
	if err != nil {
		return nil, fmt.Errorf("reading Key Vault (%s): %+v", id, err)
	}

	return utils.Bool(keyvault.FindKeyVaultAccessPolicy(resp.Properties.AccessPolicies, objectId, applicationId) != nil), nil
}

func (r KeyVaultAccessPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
  ]

  secret_permissions = [
    "Get",
    "Set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, template)
}

func (r KeyVaultAccessPolicyResource) basicMixedCase(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Get",
  ]

  secret_permissions = [
    "Get",
    "set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, template)
}

func (r KeyVaultAccessPolicyResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "import" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault_access_policy.test.tenant_id
  object_id    = azurerm_key_vault_access_policy.test.object_id

  key_permissions = [
    "Get",
  ]

  secret_permissions = [
    "Get",
    "Set",
  ]
}
`, template)
}

func (r KeyVaultAccessPolicyResource) multiple(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_access_policy" "test_with_application_id" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Create",
    "Get",
  ]

  secret_permissions = [
    "Get",
    "Delete",
  ]

  certificate_permissions = [
    "Create",
    "Delete",
  ]

  application_id = data.azurerm_client_config.current.client_id
  tenant_id      = data.azurerm_client_config.current.tenant_id
  object_id      = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_access_policy" "test_no_application_id" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "List",
    "Encrypt",
  ]

  secret_permissions = [
    "List",
    "Delete",
  ]

  certificate_permissions = [
    "List",
    "Delete",
  ]

  storage_permissions = [
    "Backup",
    "Delete",
    "DeleteSAS",
    "Get",
    "GetSAS",
    "List",
    "ListSAS",
    "Purge",
    "Recover",
    "RegenerateKey",
    "Restore",
    "Set",
    "SetSAS",
    "Update",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, template)
}

func (r KeyVaultAccessPolicyResource) nonExistentVault(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_access_policy" "test" {
  # Must appear to be URL, but not actually exist - appending a string works
  key_vault_id = "${azurerm_key_vault.test.id}NOPE"

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
  ]

  secret_permissions = [
    "Get",
  ]
}
`, template)
}

func (r KeyVaultAccessPolicyResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "List",
    "Encrypt",
  ]

  secret_permissions = []

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, template)
}

func (KeyVaultAccessPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
