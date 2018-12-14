package azurerm

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIpRanges_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureIpRangesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccAzureIpRangesCheckAttributes("data.azurerm_ip_ranges.test"),
				),
			},
		},
	})
}

func testAccAzureIpRangesCheckAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find regions data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("regions data source ID not set.")
		}

		count, ok := rs.Primary.Attributes["subnets.#"]
		if !ok {
			return errors.New("can't find 'subnets' attribute")
		}

		noOfSubnets, err := strconv.Atoi(count)
		if err != nil {
			return errors.New("failed to read number of subnets")
		}
		if noOfSubnets < 10 {
			return fmt.Errorf("expected at least 10 subnets, received %d, this is most likely a bug", noOfSubnets)
		}

		return nil
	}
}

const testAccAzureIpRangesConfig = `
data "azurerm_ip_ranges" "test" {
	regions = ["australiaeast", "australiasouth"]
}
`
