package datashare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataShareDataSetDataLakeGen1Resource struct {
}

func TestAccDataShareDataSetDataLakeGen1_basicFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen1", "test")
	r := DataShareDataSetDataLakeGen1Resource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicFile(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShareDataSetDataLakeGen1_basicFolder(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen1", "test")
	r := DataShareDataSetDataLakeGen1Resource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicFolder(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShareDataSetDataLakeGen1_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen1", "test")
	r := DataShareDataSetDataLakeGen1Resource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicFile(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t DataShareDataSetDataLakeGen1Resource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DataSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataShare.DataSetClient.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Data Share Data Set %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	switch resp := resp.Value.(type) {
	case datashare.ADLSGen1FileDataSet:
		return utils.Bool(resp.ADLSGen1FileProperties != nil), nil

	case datashare.ADLSGen1FolderDataSet:
		return utils.Bool(resp.ADLSGen1FolderProperties != nil), nil
	}

	return nil, fmt.Errorf("Data Share Data %q (Resource Group %q / accountName %q / shareName %q) is not a datalake store gen1 dataset", id.Name, id.ResourceGroup, id.AccountName, id.ShareName)
}

func (DataShareDataSetDataLakeGen1Resource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datashare-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_share_account" "test" {
  name                = "acctest-DSA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }

  tags = {
    env = "Test"
  }
}

resource "azurerm_data_share" "test" {
  name       = "acctest_DS_%[1]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"
}

resource "azurerm_data_lake_store" "test" {
  name                = "acctestdls%[3]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  firewall_state      = "Disabled"
}

resource "azurerm_data_lake_store_file" "test" {
  account_name     = azurerm_data_lake_store.test.name
  local_file_path  = "./testdata/application_gateway_test.cer"
  remote_file_path = "/test/application_gateway_test.cer"
}

data "azuread_service_principal" "test" {
  display_name = azurerm_data_share_account.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_data_lake_store.test.id
  role_definition_name = "Owner"
  principal_id         = data.azuread_service_principal.test.object_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
}

func (r DataShareDataSetDataLakeGen1Resource) basicFile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen1" "test" {
  name               = "acctest-DSDL1-%d"
  data_share_id      = azurerm_data_share.test.id
  data_lake_store_id = azurerm_data_lake_store.test.id
  file_name          = "application_gateway_test.cer"
  folder_path        = "test"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r DataShareDataSetDataLakeGen1Resource) basicFolder(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen1" "test" {
  name               = "acctest-DSDL1-%d"
  data_share_id      = azurerm_data_share.test.id
  data_lake_store_id = azurerm_data_lake_store.test.id
  folder_path        = "test"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r DataShareDataSetDataLakeGen1Resource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_share_dataset_data_lake_gen1" "import" {
  name               = azurerm_data_share_dataset_data_lake_gen1.test.name
  data_share_id      = azurerm_data_share.test.id
  data_lake_store_id = azurerm_data_share_dataset_data_lake_gen1.test.data_lake_store_id
  folder_path        = azurerm_data_share_dataset_data_lake_gen1.test.folder_path
}
`, r.basicFile(data))
}
