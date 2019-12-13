package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testCheckAzureRMSpringCloudAppExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Spring Cloud App not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group"]
		serviceName := rs.Primary.Attributes["service_name"]

		client := testAccProvider.Meta().(*ArmClient).AppPlatform.AppsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Spring Cloud App %q (Service Name %q / Resource Group %q) does not exist", name, serviceName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on appsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSpringCloudAppDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).AppPlatform.AppsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_spring_cloud_app" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group"]
		serviceName := rs.Primary.Attributes["service_name"]

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on appsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}
