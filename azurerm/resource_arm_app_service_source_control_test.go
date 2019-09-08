package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceSourceControl_ExternalGit(t *testing.T) {
	resourceName := "azurerm_app_service_source_control.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	repoUrl := "https://github.com/Azure-Samples/app-service-web-html-get-started"

	config := testAccAzureRMAppServiceSourceControlExternalGit(ri, location, repoUrl)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAppServiceSourceControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "repo_url", repoUrl),
					resource.TestCheckResourceAttr(resourceName, "type", "ExternalGit"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_service_id"},
			},
		},
	})
}

func testAccAzureRMAppServiceSourceControlExternalGit(rInt int, location, repoUrl string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%d"
  location = "%s"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctest%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctest%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  lifecycle {
    ignore_changes = [site_config.0.scm_type]
  }
}

resource "azurerm_app_service_source_control" "test" {
  app_service_id = "${azurerm_app_service.test.id}"
  type           = "ExternalGit"
  repo_url       = "%s"
  branch         = "master"
}
`, rInt, location, rInt, rInt, repoUrl)
}

func testCheckAzureRMAppServiceSourceControlDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).web.AppServicesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_source_control" {
			continue
		}

		id, err := azure.ParseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		resourceGroup := id.ResourceGroup
		name := id.Path["sites"]

		resp, err := conn.GetSourceControl(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}
