package iothub

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMIotHubSharedAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_iothub_shared_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMIotHubSharedAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMIotHubSharedAccessPolicy_basic(data acceptance.TestData) string {
	template := testAccAzureRMIotHubSharedAccessPolicy_basic(data)

	return fmt.Sprintf(`
%s

data "azurerm_iothub_shared_access_policy" "test" {
  name                = azurerm_iothub_shared_access_policy.test.name
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
}
`, template)
}
