package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_tracing", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementSubscription_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementSubscription_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementSubscription_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementSubscription_update(data, "submitted", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "submitted"),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_tracing", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
				),
			},
			{
				Config: testAccAzureRMApiManagementSubscription_update(data, "active", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "active"),
				),
			},
			{
				Config: testAccAzureRMApiManagementSubscription_update(data, "suspended", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "suspended"),
				),
			},
			{
				Config: testAccAzureRMApiManagementSubscription_update(data, "cancelled", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "cancelled"),
				),
			},
			{
				Config: testAccAzureRMApiManagementSubscription_update(data, "active", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_tracing", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMApiManagementSubscription_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementSubscription_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "active"),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_tracing", "false"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAPIManagementSubscriptionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.SubscriptionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_subscription" {
			continue
		}

		subscriptionId := rs.Primary.Attributes["subscription_id"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, subscriptionId)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}
	return nil
}

func testCheckAzureRMAPIManagementSubscriptionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.SubscriptionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		subscriptionId := rs.Primary.Attributes["subscription_id"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, subscriptionId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Subscription %q (API Management Service %q / Resource Group %q) does not exist", subscriptionId, serviceName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on apiManagement.SubscriptionsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementSubscription_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementSubscription_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_subscription" "test" {
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  user_id             = azurerm_api_management_user.test.id
  product_id          = azurerm_api_management_product.test.id
  display_name        = "Butter Parser API Enterprise Edition"
}
`, template)
}

func testAccAzureRMApiManagementSubscription_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementSubscription_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_subscription" "import" {
  subscription_id     = azurerm_api_management_subscription.test.subscription_id
  resource_group_name = azurerm_api_management_subscription.test.resource_group_name
  api_management_name = azurerm_api_management_subscription.test.api_management_name
  user_id             = azurerm_api_management_subscription.test.user_id
  product_id          = azurerm_api_management_subscription.test.product_id
  display_name        = azurerm_api_management_subscription.test.display_name
}
`, template)
}

func testAccAzureRMApiManagementSubscription_update(data acceptance.TestData, state string, allow_tracing string) string {
	template := testAccAzureRMApiManagementSubscription_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_subscription" "test" {
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  user_id             = azurerm_api_management_user.test.id
  product_id          = azurerm_api_management_product.test.id
  display_name        = "Butter Parser API Enterprise Edition"
  state               = "%s"
  allow_tracing       = "%s"
}
`, template, state, allow_tracing)
}

func testAccAzureRMApiManagementSubscription_complete(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementSubscription_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_subscription" "test" {
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  user_id             = azurerm_api_management_user.test.id
  product_id          = azurerm_api_management_product.test.id
  display_name        = "Butter Parser API Enterprise Edition"
  state               = "active"
  allow_tracing       = "false"
}
`, template)
}

func testAccAzureRMApiManagementSubscription_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = false
  published             = true
}

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
