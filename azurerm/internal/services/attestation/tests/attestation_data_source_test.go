package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAttestation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_attestation", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAttestation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
				),
			},
		},
	})
}

func testAccDataSourceAttestation_basic(data acceptance.TestData) string {
	config := testAccAzureRMAttestation_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_attestation" "test" {
  name = azurerm_attestation.test.name
  resource_group_name = azurerm_attestation.test.resource_group_name
}
`, config)
}
