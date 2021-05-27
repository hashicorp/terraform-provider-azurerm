package datashare_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DataShareDatasetDataLakeGen1DataSource struct {
}

func TestAccDataShareDatasetDataLakeGen1DataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share_dataset_data_lake_gen1", "test")
	r := DataShareDatasetDataLakeGen1DataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("data_lake_store_id").Exists(),
				check.That(data.ResourceName).Key("file_name").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
	})
}

func (DataShareDatasetDataLakeGen1DataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_share_dataset_data_lake_gen1" "test" {
  name          = azurerm_data_share_dataset_data_lake_gen1.test.name
  data_share_id = azurerm_data_share_dataset_data_lake_gen1.test.data_share_id
}
`, DataShareDataSetDataLakeGen1Resource{}.basicFile(data))
}
