package datashare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDataShareDatasetDataLakeGen1_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share_dataset_data_lake_gen1", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataShareDataSetDestroy("azurerm_data_share_dataset_data_lake_gen1"),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataShareDatasetDataLakeGen1_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "data_lake_store_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "file_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
		},
	})
}

func testAccDataSourceDataShareDatasetDataLakeGen1_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataShareDataSetDataLakeGen1File_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_share_dataset_data_lake_gen1" "test" {
  name          = azurerm_data_share_dataset_data_lake_gen1.test.name
  data_share_id = azurerm_data_share_dataset_data_lake_gen1.test.data_share_id
}
`, config)
}
