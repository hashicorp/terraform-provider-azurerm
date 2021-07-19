package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultNetworkAclRuleResource struct {
}

func TestAccKeyVaultNetworkAclRule_ip(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultNetworkAclRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ip(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			Check:             check.That(data.ResourceName).Key("network_acls.0.ip_rules.0").HasValue("1.2.3.4"),
		},
	})
}

func TestAccKeyVaultNetworkAclRule_requiresImportIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultNetworkAclRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ip(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		{
			Config:      r.requiresImportIp(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_network_acl_rule"),
		},
	})
}

func TestAccKeyVaultNetworkAclRule_cidr(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultNetworkAclRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cidr(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			Check:             check.That(data.ResourceName).Key("network_acls.0.ip_rules.1").HasValue("20.42.5.0/24"),
		},
	})
}

func TestAccKeyVaultNetworkAclRule_requiresImportCidr(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultNetworkAclRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cidr(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		{
			Config:      r.requiresImportCidr(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_network_acl_rule"),
		},
	})
}

func TestAccKeyVaultNetworkAclRule_subnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultNetworkAclRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subnet(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			Check:             check.That(data.ResourceName).Key("network_acls.0.virtual_network_subnet_ids.%").HasValue("1"),
		},
	})
}

func TestAccKeyVaultNetworkAclRule_requiresImportSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultNetworkAclRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.subnet(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		{
			Config:      r.requiresImportSubnet(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_network_acl_rule"),
		},
	})
}

func TestAccKeyVaultNetworkAclRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultNetworkAclRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("network_acls.0.ip_rules.0").HasValue("1.2.3.4"),
				check.That(data.ResourceName).Key("network_acls.0.ip_rules.1").HasValue("20.42.5.0/24"),
				check.That(data.ResourceName).Key("network_acls.0.virtual_network_subnet_ids.%").HasValue("2"),
			),
		},
	})
}

func TestAccKeyVaultNetworkAclRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	r := KeyVaultNetworkAclRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		{
			Config: r.update(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("network_acls.0.ip_rules.0").HasValue("2.3.4.5"),
				check.That(data.ResourceName).Key("network_acls.0.ip_rules.1").HasValue("40.82.252.0/24"),
				check.That(data.ResourceName).Key("network_acls.0.virtual_network_subnet_ids.%").HasValue("1"),
			),
		},
	})
}

func (KeyVaultNetworkAclRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VaultID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.KeyVault.VaultsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Key Vault (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (KeyVaultNetworkAclRuleResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VaultID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err := client.KeyVault.VaultsClient.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r KeyVaultNetworkAclRuleResource) vaultTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r KeyVaultNetworkAclRuleResource) basic(data acceptance.TestData) string {
	template := r.vaultTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault" "test" {
  name                       = "vault%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  network_acls {
    default_action = "Deny"
    bypass         = "None"
  }
}
`, template, data.RandomInteger)
}

func (r KeyVaultNetworkAclRuleResource) ip(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_network_acl_rule" "test_ip" {
  key_vault_id = azurerm_key_vault.test.id
  source       = "1.2.3.4"
}
`, template, data.RandomInteger)
}

func (r KeyVaultNetworkAclRuleResource) requiresImportIp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_network_acl_rule" "import" {
  key_vault_id = azurerm_key_vault_network_acl_rule.test_ip.key_vault_id
  source       = azurerm_key_vault_network_acl_rule.test_ip.source
}
`, r.ip(data))
}

func (r KeyVaultNetworkAclRuleResource) cidr(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_network_acl_rule" "test_cidr" {
  key_vault_id = azurerm_key_vault.test.id
  source       = "20.42.5.0/24"
}
`, template, data.RandomInteger)
}

func (r KeyVaultNetworkAclRuleResource) requiresImportCidr(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_network_acl_rule" "import" {
  key_vault_id = azurerm_key_vault_network_acl_rule.test_cidr.key_vault_id
  source       = azurerm_key_vault_network_acl_rule.test_cidr.source
}
`, r.cidr(data))
}

func (r KeyVaultNetworkAclRuleResource) networkTemplate(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test_a" {
  name                 = "acctestsubneta%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.KeyVault"]
}

resource "azurerm_subnet" "test_b" {
  name                 = "acctestsubnetb%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.4.0/24"
  service_endpoints    = ["Microsoft.KeyVault"]
}
`, template, data.RandomInteger)
}

func (r KeyVaultNetworkAclRuleResource) subnet(data acceptance.TestData) string {
	template := r.networkTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_network_acl_rule" "test_subnet" {
  key_vault_id = azurerm_key_vault.test.id
  source       = azurerm_subnet.test_a.id
}
`, template)
}

func (r KeyVaultNetworkAclRuleResource) requiresImportSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_network_acl_rule" "import" {
  key_vault_id = azurerm_key_vault_network_acl_rule.test_subnet.key_vault_id
  source       = azurerm_key_vault_network_acl_rule.test_subnet.source
}
`, r.subnet(data))
}

func (r KeyVaultNetworkAclRuleResource) complete(data acceptance.TestData) string {
	template := r.subnet(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_network_acl_rule" "test_subnet_b" {
  key_vault_id = azurerm_key_vault.test.id
  source       = azurerm_subnet.test_b.id
}

resource "azurerm_key_vault_network_acl_rule" "test_ip" {
  key_vault_id = azurerm_key_vault.test.id
  source       = "1.2.3.4"
}

resource "azurerm_key_vault_network_acl_rule" "test_cidr" {
  key_vault_id = azurerm_key_vault.test.id
  source       = "20.42.5.0/24"
}
`, template)
}

func (r KeyVaultNetworkAclRuleResource) update(data acceptance.TestData) string {
	template := r.networkTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_network_acl_rule" "test_subnet" {
  key_vault_id = azurerm_key_vault.test.id
  source       = azurerm_subnet.test_b.id
}

resource "azurerm_key_vault_network_acl_rule" "test_ip" {
  key_vault_id = azurerm_key_vault.test.id
  source       = "2.3.4.5"
}

resource "azurerm_key_vault_network_acl_rule" "test_cidr" {
  key_vault_id = azurerm_key_vault.test.id
  source       = "40.82.252.0/24"
}
`, template)
}
