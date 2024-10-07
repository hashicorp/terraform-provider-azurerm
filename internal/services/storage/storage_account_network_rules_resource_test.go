// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

func TestAccStorageAccountNetworkRules_id(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")
	r := StorageAccountNetworkRulesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.id(data),
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

func TestAccStorageAccountNetworkRules_synapseAccess(t *testing.T) {
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

func TestAccStorageAccountNetworkRules_redeploy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_network_rules", "test")
	parent := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	r := StorageAccountNetworkRulesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deploy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(parent.ResourceName).ExistsInAzure(r),
			),
		},
		parent.ImportStep(),
		{
			Config: r.remove(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(parent.ResourceName).DoesNotExistInAzure(r),
			),
		},
		parent.ImportStep(),
		{
			Config: r.deploy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(parent.ResourceName).ExistsInAzure(r),
			),
		},
		parent.ImportStep(),
	})
}

func (r StorageAccountNetworkRulesResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseStorageAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Storage.ResourceManager.StorageAccounts.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if acls := props.NetworkAcls; acls != nil {
				hasIPRules := acls.IPRules != nil && len(*acls.IPRules) > 0
				usesNonDefaultAction := acls.DefaultAction != storageaccounts.DefaultActionAllow
				usesNonDefaultBypass := acls.Bypass != nil && *acls.Bypass != storageaccounts.BypassAzureServices
				hasVirtualNetworkRules := acls.VirtualNetworkRules != nil && len(*acls.VirtualNetworkRules) > 0
				if hasIPRules || usesNonDefaultAction || usesNonDefaultBypass || hasVirtualNetworkRules {
					return pointer.To(true), nil
				}
			}
		}
	}

	return utils.Bool(false), nil
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
  address_prefixes     = ["10.0.2.0/24"]
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
  storage_account_id = azurerm_storage_account.test.id

  default_action             = "Deny"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}

func (r StorageAccountNetworkRulesResource) id(data acceptance.TestData) string {
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
  address_prefixes     = ["10.0.2.0/24"]
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
  storage_account_id = azurerm_storage_account.test.id

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
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestsubnet2%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
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
  storage_account_id = azurerm_storage_account.test.id

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
  storage_account_id = azurerm_storage_account.test.id

  default_action = "Deny"
  bypass         = ["None"]
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
  storage_account_id = azurerm_storage_account.test.id

  default_action             = "Deny"
  bypass                     = ["None"]
  ip_rules                   = []
  virtual_network_subnet_ids = []
}
`, StorageAccountResource{}.networkRulesTemplate(data), data.RandomString)
}

func (r StorageAccountNetworkRulesResource) privateLinkAccess(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_search_service" "test" {
  name                = "acctestsearchservice%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "basic"
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
  storage_account_id         = azurerm_storage_account.test.id
  default_action             = "Deny"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
  private_link_access {
    endpoint_resource_id = azurerm_search_service.test.id
  }
}
`, StorageAccountResource{}.networkRulesTemplate(data), data.RandomInteger, data.RandomString)
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

  identity {
    type = "SystemAssigned"
  }
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
  storage_account_id = azurerm_storage_account.test.id

  default_action = "Deny"
  ip_rules       = ["127.0.0.1"]
  private_link_access {
    endpoint_resource_id = azurerm_synapse_workspace.test.id
  }
}
`, StorageAccountResource{}.networkRulesTemplate(data), data.RandomString, data.RandomInteger)
}

func (r StorageAccountNetworkRulesResource) deploy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account_network_rules" "test" {
  storage_account_id = azurerm_storage_account.test.id

  default_action = "Deny"
  ip_rules       = ["198.1.1.0"]
  bypass         = ["Metrics"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageAccountNetworkRulesResource) remove(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
