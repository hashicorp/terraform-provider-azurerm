package keyvault_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultAccessPolicyResource struct {
}

func TestAccKeyVaultAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("secret_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("secret_permissions.1").HasValue("set"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKeyVaultAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("secret_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("secret_permissions.1").HasValue("set"),
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("create"),
				check.That(data.ResourceName).Key("key_permissions.1").HasValue("get"),
				check.That(data.ResourceName).Key("secret_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("secret_permissions.1").HasValue("delete"),
				check.That(data.ResourceName).Key("certificate_permissions.0").HasValue("create"),
				check.That(data.ResourceName).Key("certificate_permissions.1").HasValue("delete"),
				resource.TestCheckResourceAttr(resourceName2, "key_permissions.0", "list"),
				resource.TestCheckResourceAttr(resourceName2, "key_permissions.1", "encrypt"),
				resource.TestCheckResourceAttr(resourceName2, "secret_permissions.0", "list"),
				resource.TestCheckResourceAttr(resourceName2, "secret_permissions.1", "delete"),
				resource.TestCheckResourceAttr(resourceName2, "certificate_permissions.0", "list"),
				resource.TestCheckResourceAttr(resourceName2, "certificate_permissions.1", "delete"),
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("secret_permissions.0").HasValue("get"),
				check.That(data.ResourceName).Key("secret_permissions.1").HasValue("set"),
			),
		},
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_permissions.0").HasValue("list"),
				check.That(data.ResourceName).Key("key_permissions.1").HasValue("encrypt"),
			),
		},
	})
}

func TestAccKeyVaultAccessPolicy_nonExistentVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:             r.nonExistentVault(data),
			ExpectNonEmptyPlan: true,
			ExpectError:        regexp.MustCompile(`Error retrieving Key Vault`),
		},
	})
}

func (t KeyVaultAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
    "set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, r.template(data))
}

func (r KeyVaultAccessPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "import" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_key_vault_access_policy.test.tenant_id
  object_id    = azurerm_key_vault_access_policy.test.object_id

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
    "set",
  ]
}
`, r.basic(data))
}

func (r KeyVaultAccessPolicyResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test_with_application_id" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "create",
    "get",
  ]

  secret_permissions = [
    "get",
    "delete",
  ]

  certificate_permissions = [
    "create",
    "delete",
  ]

  application_id = data.azurerm_client_config.current.client_id
  tenant_id      = data.azurerm_client_config.current.tenant_id
  object_id      = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_access_policy" "test_no_application_id" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "list",
    "encrypt",
  ]

  secret_permissions = [
    "list",
    "delete",
  ]

  certificate_permissions = [
    "list",
    "delete",
  ]

  storage_permissions = [
    "backup",
    "delete",
    "deletesas",
    "get",
    "getsas",
    "list",
    "listsas",
    "purge",
    "recover",
    "regeneratekey",
    "restore",
    "set",
    "setsas",
    "update",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, r.template(data))
}

func (r KeyVaultAccessPolicyResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "list",
    "encrypt",
  ]

  secret_permissions = []

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}
`, r.template(data))
}

func (KeyVaultAccessPolicyResource) template(data acceptance.TestData) string {
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
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "premium"

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (KeyVaultAccessPolicyResource) nonExistentVault(data acceptance.TestData) string {
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
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  # Must appear to be URL, but not actually exist - appending a string works
  key_vault_id = "${azurerm_key_vault.test.id}NOPE"

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
