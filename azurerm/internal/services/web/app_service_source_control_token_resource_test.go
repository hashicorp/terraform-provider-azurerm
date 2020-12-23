package web_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAppServiceSourceControlToken(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_source_control_token", "test")
	token := strings.ToLower(acctest.RandString(41))
	tokenSecret := strings.ToLower(acctest.RandString(41))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAppServiceSourceControlTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceSourceControlToken(token, tokenSecret),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "type", "GitHub"),
					resource.TestCheckResourceAttr(data.ResourceName, "token", token),
					resource.TestCheckResourceAttr(data.ResourceName, "token_secret", tokenSecret),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAppServiceSourceControlToken(token, tokenSecret string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_app_service_source_control_token" "test" {
  type         = "GitHub"
  token        = "%s"
  token_secret = "%s"
}
`, token, tokenSecret)
}

func testCheckAppServiceSourceControlTokenDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Web.BaseClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_source_control_token" {
			continue
		}

		scmType := rs.Primary.Attributes["type"]

		resp, err := conn.GetSourceControl(ctx, scmType)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}
