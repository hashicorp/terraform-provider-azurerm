package azurerm

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform/helper/resource"
    "github.com/hashicorp/terraform/terraform"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)


func testCheckAzureRMFrontDoorExists(resourceName string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok {
            return fmt.Errorf("Front Door not found: %s", resourceName)
        }

        name := rs.Primary.Attributes["name"]
        resourceGroup := rs.Primary.Attributes["resource_group_name"]

        client := testAccProvider.Meta().(*ArmClient).frontdoorClient
        ctx := testAccProvider.Meta().(*ArmClient).StopContext

        if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
            if utils.ResponseWasNotFound(resp.Response) {
                return fmt.Errorf("Bad: Front Door %q (Resource Group %q) does not exist", name, resourceGroup)
            }
            return fmt.Errorf("Bad: Get on frontdoorClient: %+v", err)
        }

        return nil
    }
}

func testCheckAzureRMFrontDoorDestroy(s *terraform.State) error {
    client := testAccProvider.Meta().(*ArmClient).frontdoorClient
    ctx := testAccProvider.Meta().(*ArmClient).StopContext

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "azurerm_front_door" {
            continue
        }

        name := rs.Primary.Attributes["name"]
        resourceGroup := rs.Primary.Attributes["resource_group_name"]

        if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
            if !utils.ResponseWasNotFound(resp.Response) {
                return fmt.Errorf("Bad: Get on frontdoorClient: %+v", err)
            }
        }

        return nil
    }

    return nil
}
