package authorization

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMClientConfig_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_client_config", "current")
	clientId := os.Getenv("ARM_CLIENT_ID")
	tenantId := os.Getenv("ARM_TENANT_ID")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckArmClientConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "client_id", clientId),
					resource.TestCheckResourceAttr(data.ResourceName, "tenant_id", tenantId),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_id", subscriptionId),
					testAzureRMClientConfigGUIDAttr(data.ResourceName, "object_id"),
				),
			},
		},
	})
}

func testAzureRMClientConfigGUIDAttr(name, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := regexp.MustCompile("^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$")

		return resource.TestMatchResourceAttr(name, key, r)(s)
	}
}

const testAccCheckArmClientConfig_basic = `
data "azurerm_client_config" "current" { }
`
