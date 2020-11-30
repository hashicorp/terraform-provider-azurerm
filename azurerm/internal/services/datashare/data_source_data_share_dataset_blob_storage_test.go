package datashare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceDataShareDatasetBlobStorage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share_dataset_blob_storage", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareDataSetDestroy("azurerm_data_share_dataset_blob_storage"),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataShareDatasetBlobStorage_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account.0.name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account.0.resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account.0.subscription_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "file_path"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
		},
	})
}

func testAccDataSourceDataShareDatasetBlobStorage_basic(data acceptance.TestData) string {
	config := testAccDataShareDataSetBlobStorageFile_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_share_dataset_blob_storage" "test" {
  name          = azurerm_data_share_dataset_blob_storage.test.name
  data_share_id = azurerm_data_share_dataset_blob_storage.test.data_share_id
}
`, config)
}
