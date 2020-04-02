package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourcePeerAsn_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_peer_asn", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePeerAsn_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourcePeerAsn_basic(data acceptance.TestData) string {
	config := testAccAzureRMPeerAsn_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_peer_asn" "test" {
  name = azurerm_peer_asn.test.name
}
`, config)
}
