package datashare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceDataShareKustoClusterDataset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share_dataset_kusto_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataShareDataSetDestroy("azurerm_data_share_dataset_kusto_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataShareKustoClusterDataset_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareDataSetExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kusto_cluster_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kusto_cluster_location"),
				),
			},
		},
	})
}

func testAccDataSourceDataShareKustoClusterDataset_basic(data acceptance.TestData) string {
	config := testAccDataShareKustoClusterDataSet_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_share_dataset_kusto_cluster" "test" {
  name     = azurerm_data_share_dataset_kusto_cluster.test.name
  share_id = azurerm_data_share_dataset_kusto_cluster.test.share_id
}
`, config)
}
