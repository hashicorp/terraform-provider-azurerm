package databox_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataBoxJobDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_box_job", "test")

	resource.ParallelTest(t, resource.TestCase{
		// nolint missing CheckDestroy
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataBoxJob_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceDataBoxJob_basic(data acceptance.TestData) string {
	config := testAccDataBoxJob_existingResource(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_box_job" "test" {
  name                = azurerm_data_box_job.test.name
  resource_group_name = azurerm_data_box_job.test.resource_group_name
}
`, config)
}

func testAccDataBoxJob_existingResource(data acceptance.TestData) string {
	template := testAccDataBoxJob_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_box_job" "test" {
  name                = "acctest-DataBox-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  contact_details {
    name         = "Terraform Test"
    emails       = ["some.user@example.com"]
    phone_number = "+11234567891"
  }

  destination_storage_account {
    storage_account_id = azurerm_storage_account.test.id
  }

  shipping_address {
    city              = "Redmond"
    country           = "US"
    postal_code       = "98052"
    state_or_province = "WA"
    street_address_1  = "One Microsoft Way"
  }

  preferred_shipment_type = "MicrosoftManaged"

  sku_name = "DataBox"
}
`, template, data.RandomString)
}
