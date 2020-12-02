package datashare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DataShareDatasetBlobStorageDataSource struct {
}

func TestAccDataShareDatasetBlobStorageDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share_dataset_blob_storage", "test")
	r := DataShareDatasetBlobStorageDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(DataShareDataSetBlobStorageResource{}),
				check.That(data.ResourceName).Key("container_name").Exists(),
				check.That(data.ResourceName).Key("storage_account.0.name").Exists(),
				check.That(data.ResourceName).Key("storage_account.0.resource_group_name").Exists(),
				check.That(data.ResourceName).Key("storage_account.0.subscription_id").Exists(),
				check.That(data.ResourceName).Key("file_path").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
	})
}

func (DataShareDatasetBlobStorageDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_share_dataset_blob_storage" "test" {
  name          = azurerm_data_share_dataset_blob_storage.test.name
  data_share_id = azurerm_data_share_dataset_blob_storage.test.data_share_id
}
`, DataShareDataSetBlobStorageResource{}.basicFile(data))
}
