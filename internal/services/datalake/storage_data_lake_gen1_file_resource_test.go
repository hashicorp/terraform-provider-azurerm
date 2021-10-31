package datalake_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageDataLakeGen1FileResource struct {
}

func TestAccStorageDataLakeGen1File_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_file", "test")
	r := StorageDataLakeGen1FileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("local_file_path"),
	})
}

func TestAccStorageDataLakeGen1File_largefiles(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_file", "test")
	r := StorageDataLakeGen1FileResource{}

	// "large" in this context is anything greater than 4 megabytes
	largeSize := 12 * 1024 * 1024 // 12 mb
	bytes := make([]byte, largeSize)
	rand.Read(bytes) // fill with random data

	tmpfile, err := os.CreateTemp("", "azurerm-acc-datalake-file-large")
	if err != nil {
		t.Errorf("Unable to open a temporary file.")
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(bytes); err != nil {
		t.Errorf("Unable to write to temporary file %q: %v", tmpfile.Name(), err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Errorf("Unable to close temporary file %q: %v", tmpfile.Name(), err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.largefiles(data, tmpfile.Name()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("local_file_path"),
	})
}

func TestAccStorageDataLakeGen1File_requiresimport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen1_file", "test")
	r := StorageDataLakeGen1FileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_storage_data_lake_gen1_file"),
		},
	})
}

func (t StorageDataLakeGen1FileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Datalake.StoreFilesClient
	id, err := datalake.ParseStorageDataLakeGen1FileId(state.ID, client.AdlsFileSystemDNSSuffix)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetFileStatus(ctx, id.StorageAccountName, id.FilePath, utils.Bool(true))
	if err != nil {
		return nil, fmt.Errorf("retrieving Storage Date Lake Gen1 File Rule %q (Account %q): %v", id.FilePath, id.StorageAccountName, err)
	}

	return utils.Bool(resp.FileStatus != nil), nil
}

func (StorageDataLakeGen1FileResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  firewall_state      = "Disabled"
}

resource "azurerm_storage_data_lake_gen1_file" "test" {
  remote_file_path = "/test/application_gateway_test.cer"
  account_name     = azurerm_storage_data_lake_gen1_filesystem.test.name
  local_file_path  = "./testdata/application_gateway_test.cer"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Primary)
}

func (StorageDataLakeGen1FileResource) largefiles(data acceptance.TestData, file string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datalake-%d"
  location = "%s"
}

resource "azurerm_storage_data_lake_gen1_filesystem" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  firewall_state      = "Disabled"
}

resource "azurerm_storage_data_lake_gen1_file" "test" {
  remote_file_path = "/test/testAccAzureRMStorageDataLakeGen1File_largefiles.bin"
  account_name     = azurerm_storage_data_lake_gen1_filesystem.test.name
  local_file_path  = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.Locations.Primary, file)
}

func (StorageDataLakeGen1FileResource) requiresImport(data acceptance.TestData) string {
	template := StorageDataLakeGen1FileResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen1_file" "import" {
  remote_file_path = azurerm_storage_data_lake_gen1_file.test.remote_file_path
  account_name     = azurerm_storage_data_lake_gen1_file.test.account_name
  local_file_path  = "./testdata/application_gateway_test.cer"
}
`, template)
}
