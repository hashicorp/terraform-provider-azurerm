package datashare_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataShareDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataShare_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "account_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kind"),
				),
			},
		},
	})
}

func TestAccDataShareDataSource_snapshotSchedule(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_share", "test")
	startTime := time.Now().Add(time.Hour * 7).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataShare_snapshotSchedule(data, startTime),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataShareExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "snapshot_schedule.0.name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "snapshot_schedule.0.recurrence"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "snapshot_schedule.0.start_time"),
				),
			},
		},
	})
}

func testAccDataSourceDataShare_basic(data acceptance.TestData) string {
	config := testAccDataShare_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_share" "test" {
  name       = azurerm_data_share.test.name
  account_id = azurerm_data_share_account.test.id
}
`, config)
}

func testAccDataSourceDataShare_snapshotSchedule(data acceptance.TestData, startTime string) string {
	config := testAccDataShare_snapshotSchedule(data, startTime)
	return fmt.Sprintf(`
%s

data "azurerm_data_share" "test" {
  name       = azurerm_data_share.test.name
  account_id = azurerm_data_share_account.test.id
}
`, config)
}
