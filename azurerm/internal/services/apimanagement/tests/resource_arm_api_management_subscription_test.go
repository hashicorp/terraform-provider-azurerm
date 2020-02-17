package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAPIManagementSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAPIManagementSubscription_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_api_management_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementSubscription_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAPIManagementSubscription_requiresImport),
		},
	})
}

func TestAccAzureRMAPIManagementSubscription_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementSubscription_update(data, "submitted"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "submitted"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
				),
			},
			{
				Config: testAccAzureRMAPIManagementSubscription_update(data, "active"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "active"),
				),
			},
			{
				Config: testAccAzureRMAPIManagementSubscription_update(data, "suspended"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "suspended"),
				),
			},
			{
				Config: testAccAzureRMAPIManagementSubscription_update(data, "cancelled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "cancelled"),
				),
			},
		},
	})
}

func TestAccAzureRMAPIManagementSubscription_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_subscription", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementSubscription_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "active"),
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

func testAccAzureRMAPIManagementSubscription_basic(data acceptance.TestData) string {
	template := testAccAzureRMAPIManagementSubscription_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_subscription" "test" {
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  user_id             = "${azurerm_api_management_user.test.id}"
  product_id          = "${azurerm_api_management_product.test.id}"
  display_name        = "Butter Parser API Enterprise Edition"
}
`, template)
}

func testAccAzureRMAPIManagementSubscription_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAPIManagementSubscription_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_subscription" "import" {
  resource_group_name = "${azurerm_api_management_subscription.test.resource_group_name}"
  api_management_name = "${azurerm_api_management_subscription.test.api_management_name}"
  user_id             = "${azurerm_api_management_subscription.test.user_id}"
  product_id          = "${azurerm_api_management_subscription.test.product_id}"
  display_name        = "${azurerm_api_management_subscription.test.display_name}"
}
`, template)
}

func testAccAzureRMAPIManagementSubscription_update(data acceptance.TestData, state string) string {
	template := testAccAzureRMAPIManagementSubscription_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_subscription" "test" {
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  user_id             = "${azurerm_api_management_user.test.id}"
  product_id          = "${azurerm_api_management_product.test.id}"
  display_name        = "Butter Parser API Enterprise Edition"
  state               = "%s"
}
`, template, state)
}

func testAccAzureRMAPIManagementSubscription_complete(data acceptance.TestData) string {
	template := testAccAzureRMAPIManagementSubscription_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_subscription" "test" {
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  user_id             = "${azurerm_api_management_user.test.id}"
  product_id          = "${azurerm_api_management_product.test.id}"
  display_name        = "Butter Parser API Enterprise Edition"
  state               = "active"
}
`, template)
}

func testAccAzureRMAPIManagementSubscription_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = "${azurerm_api_management.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = false
  published             = true
}

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
