package datashare_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DataShareDatasetDataLakeGen2DataSource struct {
}

func TestAccDataShareDatasetDataLakeGen2DataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share_dataset_data_lake_gen2", "test")
	r := DataShareDatasetDataLakeGen2DataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("file_system_name").Exists(),
				check.That(data.ResourceName).Key("file_path").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
	})
}

func (DataShareDatasetDataLakeGen2DataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_share_dataset_data_lake_gen2" "test" {
  name     = azurerm_data_share_dataset_data_lake_gen2.test.name
  share_id = azurerm_data_share_dataset_data_lake_gen2.test.share_id
}
`, DataShareDataSetDataLakeGen2Resource{}.basicFile(data))
}
