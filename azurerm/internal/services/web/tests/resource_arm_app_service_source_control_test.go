package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceSourceControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control", "test")
	repoUrl := "https://github.com/Azure-Samples/app-service-web-html-get-started"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceSourceControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceSourceControl(data, repoUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "repo_url", repoUrl),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_service_id"},
			},
		},
	})
}

func testAccAzureRMAppServiceSourceControl(data acceptance.TestData, repoUrl string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appservice-%d"
  location = "%s"
}
resource "azurerm_app_service_plan" "test" {
  name                = "acctest-ASP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku {
    tier = "Standard"
    size = "S1"
  }
}
resource "azurerm_app_service" "test" {
  name                = "acctest-AS-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
  lifecycle {
    ignore_changes = [site_config.0.scm_type]
  }
}
resource "azurerm_app_service_source_control" "test" {
  app_service_id             = "${azurerm_app_service.test.id}"
  repo_url                   = "%s"
  branch                     = "master"
  is_manual_integration      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, repoUrl)
}

func testCheckAzureRMAppServiceSourceControlDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_source_control" {
			continue
		}

		id := rs.Primary.Attributes["id"]
		parsedID, err := azure.ParseAzureResourceID(id)
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}
		name := parsedID.Path["sites"]
		resourceGroup := parsedID.ResourceGroup

		resp, err := client.GetSourceControl(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}
