package datashare_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DataShareDataSource struct {
}

func TestAccDataShareDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share", "test")
	r := DataShareDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("account_id").Exists(),
				check.That(data.ResourceName).Key("kind").Exists(),
			),
		},
	})
}

func TestAccDataShareDataSource_snapshotSchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share", "test")
	r := DataShareDataSource{}
	startTime := time.Now().Add(time.Hour * 7).Format(time.RFC3339)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.snapshotSchedule(data, startTime),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("snapshot_schedule.0.name").Exists(),
				check.That(data.ResourceName).Key("snapshot_schedule.0.recurrence").Exists(),
				check.That(data.ResourceName).Key("snapshot_schedule.0.start_time").Exists(),
			),
		},
	})
}

func (DataShareDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_share" "test" {
  name       = azurerm_data_share.test.name
  account_id = azurerm_data_share_account.test.id
}
`, DataShareResource{}.basic(data))
}

func (DataShareDataSource) snapshotSchedule(data acceptance.TestData, startTime string) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_share" "test" {
  name       = azurerm_data_share.test.name
  account_id = azurerm_data_share_account.test.id
}
`, DataShareResource{}.snapshotSchedule(data, startTime))
}
