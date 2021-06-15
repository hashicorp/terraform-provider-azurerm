package datashare_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DataShareDatasetKustoDatabaseDataSource struct {
}

func TestAccDataShareDatasetKustoDatabaseDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share_dataset_kusto_database", "test")
	r := DataShareDatasetKustoDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kusto_database_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("kusto_cluster_location").Exists(),
			),
		},
	})
}

func (DataShareDatasetKustoDatabaseDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_share_dataset_kusto_database" "test" {
  name     = azurerm_data_share_dataset_kusto_database.test.name
  share_id = azurerm_data_share_dataset_kusto_database.test.share_id
}
`, DataShareDataSetKustoDatabaseResource{}.basic(data))
}
