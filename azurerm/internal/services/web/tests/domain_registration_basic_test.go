package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDomainRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_domain_registration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		CheckDestroy: checkDomainRegistrationIsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testDomainRegistration_basic(data),
				Check: resource.ComposeTestCheckFunc(
					checkDomainRegistrationExists(data.ResourceName),
					),
			},
			data.ImportStep(),
		},
	})
}

func testDomainRegistration_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_domain_registration" "test" {
  name = "acctest"
}
`)
}