package storage_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageAccountNetworkRuleResource struct{}

func TestAccStorageAccountNetworkRule_ipRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rule", "test")
	r := StorageAccountNetworkRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_rule").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("virtual_network_rule").DoesNotExist(),
				check.That(data.ResourceName).Key("resource_access_rule").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountNetworkRule_virtualNetworkRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rule", "test")
	r := StorageAccountNetworkRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetworkRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_rule").DoesNotExist(),
				check.That(data.ResourceName).Key("virtual_network_rule").Exists(),
				check.That(data.ResourceName).Key("resource_access_rule").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountNetworkRule_resourceAccessRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rule", "test")
	r := StorageAccountNetworkRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.resourceAccessRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_rule").DoesNotExist(),
				check.That(data.ResourceName).Key("virtual_network_rule").DoesNotExist(),
				check.That(data.ResourceName).Key("resource_access_rule.0.resource_id").Exists(),
				check.That(data.ResourceName).Key("resource_access_rule.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageAccountNetworkRuleResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.StorageAccountNetworkRuleID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %+v", state.ID, err)
	}

	resp, err := client.Storage.AccountsClient.GetProperties(ctx, id.StorageAccountId.ResourceGroup, id.StorageAccountId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id.StorageAccountId.ID(), err)
	}

	if rules := resp.NetworkRuleSet; rules != nil {
		if _, _, exists := storage.FindStorageAccountNetworkIPRule(rules.IPRules, id.IPRule); exists {
			return utils.Bool(true), nil
		}
		if _, _, exists := storage.FindStorageAccountVirtualNetworkRule(rules.VirtualNetworkRules, id.VirtualNetworkRule); exists {
			return utils.Bool(true), nil
		}
		if _, _, exists := storage.FindStorageAccountNetworkResourceAccessRule(rules.ResourceAccessRules, id.ResourceAccessRule); exists {
			return utils.Bool(true), nil
		}
	}
	return utils.Bool(false), nil
}

func (r StorageAccountNetworkRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r StorageAccountNetworkRuleResource) ipRule(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_network_rule" "test" {
  storage_account_id = azurerm_storage_account.test.id
  ip_rule            = "127.0.0.1"
}
`, template)
}

func (r StorageAccountNetworkRuleResource) virtualNetworkRule(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account_network_rule" "test" {
  storage_account_id   = azurerm_storage_account.test.id
  virtual_network_rule = azurerm_subnet.test.id
}
`, template)
}

func (r StorageAccountNetworkRuleResource) resourceAccessRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_account_network_rule" "test" {
  storage_account_id = azurerm_storage_account.test.id

  resource_access_rule {
    resource_id = azurerm_private_endpoint.blob.id
  }
}
`, StorageAccountResource{}.networkRulesPrivateEndpointTemplate(data), data.RandomString)
}
