package azurerm

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMClientConfig_basic(t *testing.T) {
	dataSourceName := "data.azurerm_client_config.current"
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
					resource.TestCheckResourceAttr(dataSourceName, "client_id", clientId),
					resource.TestCheckResourceAttr(dataSourceName, "tenant_id", tenantId),
					resource.TestCheckResourceAttr(dataSourceName, "subscription_id", subscriptionId),
					testAzureRMClientConfigGUIDAttr(dataSourceName, "object_id"),
					testAzureRMClientConfigGUIDAttr(dataSourceName, "service_principal_application_id"),
					testAzureRMClientConfigGUIDAttr(dataSourceName, "service_principal_object_id"),
				),
			},
		},
	})
}

func testAzureRMClientConfigGUIDAttr(name, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, err := regexp.Compile("^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$")
		if err != nil {
			return err
		}

		return resource.TestMatchResourceAttr(name, key, r)(s)
	}
}

const testAccCheckArmClientConfig_basic = `
data "azurerm_client_config" "current" { }
`
