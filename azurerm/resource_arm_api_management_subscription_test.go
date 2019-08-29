package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAPIManagementSubscription_basic(t *testing.T) {
	resourceName := "azurerm_api_management_subscription.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementSubscription_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAPIManagementSubscription_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_subscription.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementSubscription_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
				),
			},
			{
				Config:      testAccAzureRMAPIManagementSubscription_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_subscription"),
			},
		},
	})
}

func TestAccAzureRMAPIManagementSubscription_update(t *testing.T) {
	resourceName := "azurerm_api_management_subscription.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementSubscription_update(ri, location, "submitted"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "submitted"),
					resource.TestCheckResourceAttrSet(resourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
				),
			},
			{
				Config: testAccAzureRMAPIManagementSubscription_update(ri, location, "active"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "active"),
				),
			},
			{
				Config: testAccAzureRMAPIManagementSubscription_update(ri, location, "suspended"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "suspended"),
				),
			},
			{
				Config: testAccAzureRMAPIManagementSubscription_update(ri, location, "cancelled"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "cancelled"),
				),
			},
		},
	})
}

func TestAccAzureRMAPIManagementSubscription_complete(t *testing.T) {
	resourceName := "azurerm_api_management_subscription.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementSubscription_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementSubscriptionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "state", "active"),
					resource.TestCheckResourceAttrSet(resourceName, "subscription_id"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMAPIManagementSubscriptionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).apiManagement.SubscriptionsClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_subscription" {
			continue
		}

		subscriptionId := rs.Primary.Attributes["subscription_id"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		subscriptionId := rs.Primary.Attributes["subscription_id"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := testAccProvider.Meta().(*ArmClient).apiManagement.SubscriptionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testAccAzureRMAPIManagementSubscription_basic(rInt int, location string) string {
	template := testAccAzureRMAPIManagementSubscription_template(rInt, location)
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

func testAccAzureRMAPIManagementSubscription_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAPIManagementSubscription_basic(rInt, location)
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

func testAccAzureRMAPIManagementSubscription_update(rInt int, location, state string) string {
	template := testAccAzureRMAPIManagementSubscription_template(rInt, location)
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

func testAccAzureRMAPIManagementSubscription_complete(rInt int, location string) string {
	template := testAccAzureRMAPIManagementSubscription_template(rInt, location)
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

func testAccAzureRMAPIManagementSubscription_template(rInt int, location string) string {
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

  sku {
    name     = "Developer"
    capacity = 1
  }
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
`, rInt, location, rInt, rInt, rInt)
}
