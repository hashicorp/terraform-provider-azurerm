package datalake_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageDataLakeGen1FilesystemResource struct {
}

func TestAccStorageDataLakeGen1FileSystem_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_filesystem", "test")
	r := StorageDataLakeGen1FilesystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Consumption"),
				check.That(data.ResourceName).Key("encryption_state").HasValue("Enabled"),
				check.That(data.ResourceName).Key("encryption_type").HasValue("ServiceManaged"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen1FileSystem_withIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_filesystem", "test")
	r := StorageDataLakeGen1FilesystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Consumption"),
				check.That(data.ResourceName).Key("encryption_state").HasValue("Enabled"),
				check.That(data.ResourceName).Key("encryption_type").HasValue("ServiceManaged"),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen1FileSystem_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_filesystem", "test")
	r := StorageDataLakeGen1FilesystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_storage_data_lake_gen1_filesystem"),
		},
	})
}

func TestAccStorageDataLakeGen1FileSystem_tier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_filesystem", "test")
	r := StorageDataLakeGen1FilesystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tier(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tier").HasValue("Commitment_1TB"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen1FileSystem_encryptionDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_filesystem", "test")
	r := StorageDataLakeGen1FilesystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("encryption_state").HasValue("Disabled"),
				check.That(data.ResourceName).Key("encryption_type").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen1FileSystem_firewallUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_filesystem", "test")
	r := StorageDataLakeGen1FilesystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.firewall(data, "Enabled", "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("firewall_state").HasValue("Enabled"),
				check.That(data.ResourceName).Key("firewall_allow_azure_ips").HasValue("Enabled"),
			),
		},
		{
			Config: r.firewall(data, "Enabled", "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("firewall_state").HasValue("Enabled"),
				check.That(data.ResourceName).Key("firewall_allow_azure_ips").HasValue("Disabled"),
			),
		},
		{
			Config: r.firewall(data, "Disabled", "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("firewall_state").HasValue("Disabled"),
				check.That(data.ResourceName).Key("firewall_allow_azure_ips").HasValue("Enabled"),
			),
		},
		{
			Config: r.firewall(data, "Disabled", "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("firewall_state").HasValue("Disabled"),
				check.That(data.ResourceName).Key("firewall_allow_azure_ips").HasValue("Disabled"),
			),
		},
	})
}

func TestAccStorageDataLakeGen1FileSystem_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_filesystem", "test")
	r := StorageDataLakeGen1FilesystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (t StorageDataLakeGen1FilesystemResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Datalake.StoreAccountsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Storage Data Lake Gen1 Filesystem %s: %+v", id, err)
	}

	return utils.Bool(resp.DataLakeStoreAccountProperties != nil), nil
}

func (StorageDataLakeGen1FilesystemResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}

func (StorageDataLakeGen1FilesystemResource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}

func (StorageDataLakeGen1FilesystemResource) requiresImport(data acceptance.TestData) string {
	template := StorageDataLakeGen1FilesystemResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen1_filesystem" "import" {
  name                = azurerm_storage_data_lake_gen1_filesystem.test.name
  resource_group_name = azurerm_storage_data_lake_gen1_filesystem.test.resource_group_name
  location            = azurerm_storage_data_lake_gen1_filesystem.test.location
}
`, template)
}

func (StorageDataLakeGen1FilesystemResource) tier(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tier                = "Commitment_1TB"
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}

func (StorageDataLakeGen1FilesystemResource) encryptionDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  encryption_state    = "Disabled"
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}

func (StorageDataLakeGen1FilesystemResource) firewall(data acceptance.TestData, firewallState string, firewallAllowAzureIPs string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                     = "acctest%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  firewall_state           = "%s"
  firewall_allow_azure_ips = "%s"
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17], firewallState, firewallAllowAzureIPs)
}

func (StorageDataLakeGen1FilesystemResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}

func (StorageDataLakeGen1FilesystemResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, strconv.Itoa(data.RandomInteger)[2:17])
}
