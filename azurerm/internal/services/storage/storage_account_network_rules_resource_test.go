package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StorageAccountNetworkRulesResource struct{}

func TestAccStorageAccountNetworkRules_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")
	r := StorageAccountNetworkRulesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountNetworkRules_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")
	r := StorageAccountNetworkRulesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountNetworkRules_privateLinkAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")
	r := StorageAccountNetworkRulesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disablePrivateLinkAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.privateLinkAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.disablePrivateLinkAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountNetworkRules_SynapseAccess(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")
	r := StorageAccountNetworkRulesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.disablePrivateLinkAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.synapseAccess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountNetworkRules_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")
	r := StorageAccountNetworkRulesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.empty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_storage_account.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageAccountNetworkRulesResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	storageAccountName := state.Attributes["storage_account_name"]
	resourceGroup := state.Attributes["resource_group_name"]

	resp, err := client.Storage.AccountsClient.GetProperties(ctx, resourceGroup, storageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Storage Account %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r StorageAccountNetworkRulesResource) basic(data acceptance.TestData) string {
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

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Deny"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r StorageAccountNetworkRulesResource) update(data acceptance.TestData) string {
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

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.3.0/24"
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

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Allow"
  ip_rules                   = ["127.0.0.2", "127.0.0.3"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id, azurerm_subnet.test2.id]
  bypass                     = ["Metrics"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r StorageAccountNetworkRulesResource) empty(data acceptance.TestData) string {
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

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Deny"
  bypass                     = ["None"]
  ip_rules                   = []
  virtual_network_subnet_ids = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageAccountNetworkRulesResource) disablePrivateLinkAccess(data acceptance.TestData) string {
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

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Deny"
  bypass                     = ["None"]
  ip_rules                   = []
  virtual_network_subnet_ids = []
}
`, StorageAccountResource{}.networkRulesPrivateEndpointTemplate(data), data.RandomString)
}

func (r StorageAccountNetworkRulesResource) privateLinkAccess(data acceptance.TestData) string {
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

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Deny"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
  private_link_access {
    endpoint_resource_id = azurerm_private_endpoint.blob.id
  }
  private_link_access {
    endpoint_resource_id = azurerm_private_endpoint.table.id
  }
}
`, StorageAccountResource{}.networkRulesPrivateEndpointTemplate(data), data.RandomString)
}

func (r StorageAccountNetworkRulesResource) synapseAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "synapse" {
  name                     = "acctestacc%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[3]d"
  storage_account_id = azurerm_storage_account.synapse.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%[3]d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}


resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action = "Deny"
  ip_rules       = ["127.0.0.1"]
  private_link_access {
    endpoint_resource_id = azurerm_synapse_workspace.test.id
  }
}
`, StorageAccountResource{}.networkRulesPrivateEndpointTemplate(data), data.RandomString, data.RandomInteger)
}
