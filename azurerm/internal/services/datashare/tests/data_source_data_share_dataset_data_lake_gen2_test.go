package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDataShareDatasetDataLakeGen2_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share_dataset_data_lake_gen2", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_data_lake_gen2"),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataShareDatasetDataLakeGen2_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "file_system_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "file_path"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
		},
	})
}

func testAccDataSourceDataShareDatasetDataLakeGen2_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetDataLakeGen2File_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_share_dataset_data_lake_gen2" "test" {
  name     = azurerm_data_share_dataset_data_lake_gen2.test.name
  share_id = azurerm_data_share_dataset_data_lake_gen2.test.share_id
}
`, config)
}
