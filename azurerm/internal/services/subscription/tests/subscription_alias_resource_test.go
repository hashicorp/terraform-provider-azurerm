package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMsubscriptionAlias_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_alias", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMsubscriptionAliasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMsubscriptionAlias_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMsubscriptionAliasExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMsubscriptionAlias_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_alias", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMsubscriptionAliasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMsubscriptionAlias_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMsubscriptionAliasExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMsubscriptionAlias_requiresImport),
		},
	})
}

func testCheckAzureRMsubscriptionAliasExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Subscription.AliasClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("subscription Alias not found: %s", resourceName)
		}
		id, err := parse.SubscriptionAliasID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Subscription Alias %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Subscription.AliasClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMsubscriptionAliasDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Subscription.AliasClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_subscription_alias" {
			continue
		}
		id, err := parse.SubscriptionAliasID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Subscription.AliasClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMsubscriptionAlias_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_subscription_alias" "test" {
  name            = "acctest-sa-%d"
  subscription_id = data.azurerm_subscription.current.subscription_id
}
`, data.RandomInteger)
}

func testAccAzureRMsubscriptionAlias_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMsubscriptionAlias_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subscription_alias" "import" {
  name            = azurerm_subscription_alias.test.name
  subscription_id = azurerm_subscription_alias.test.subscription_id
}
`, config)
}
