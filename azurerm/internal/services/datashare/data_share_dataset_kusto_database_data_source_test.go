package datashare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataShareDatasetKustoDatabaseDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share_dataset_kusto_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareDataSetDestroy("azurerm_data_share_dataset_kusto_database"),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataShareDatasetKustoDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kusto_database_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kusto_cluster_location"),
				),
			},
		},
	})
}

func testAccDataSourceDataShareDatasetKustoDatabase_basic(data acceptance.TestData) string {
	config := testAccDataShareDataSetKustoDatabase_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_share_dataset_kusto_database" "test" {
  name     = azurerm_data_share_dataset_kusto_database.test.name
  share_id = azurerm_data_share_dataset_kusto_database.test.share_id
}
`, config)
}
