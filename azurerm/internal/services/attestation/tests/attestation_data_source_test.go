package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAttestationProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_attestation_provider", "test")
	randStr := strings.ToLower(acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAttestationProvider_basic(data, randStr),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationProviderExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccDataSourceAttestationProvider_basic(data acceptance.TestData, randStr string) string {
	config := testAccAzureRMAttestationProvider_basic(data, randStr)
	return fmt.Sprintf(`
%s

data "azurerm_attestation_provider" "test" {
  name                = azurerm_attestation_provider.test.name
  resource_group_name = azurerm_attestation_provider.test.resource_group_name
}
`, config)
}
