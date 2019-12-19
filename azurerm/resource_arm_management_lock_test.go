package azurerm

import (
	"fmt"
	"testing"

	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestValidateManagementLockName(t *testing.T) {
	str := acctest.RandString(259)
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"ab", false},
		{"ABC", false},
		{"abc", false},
		{"abc123ABC", false},
		{"123abcABC", false},
		{"ABC123abc", false},
		{"abc-123", false},
		{"abc_123", false},
		{str, false},
		{str + "h", true},
	}

	for _, test := range testCases {
		_, es := validateArmManagementLockName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

func TestAccAzureRMManagementLock_resourceGroupReadOnlyBasic(t *testing.T) {
	resourceName := "azurerm_management_lock.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMManagementLock_resourceGroupReadOnlyBasic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(resourceName),
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

func TestAccAzureRMManagementLock_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_management_lock.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMManagementLock_resourceGroupReadOnlyBasic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMManagementLock_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_management_lock"),
			},
		},
	})
}

func TestAccAzureRMManagementLock_resourceGroupReadOnlyComplete(t *testing.T) {
	resourceName := "azurerm_management_lock.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMManagementLock_resourceGroupReadOnlyComplete(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(resourceName),
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

func TestAccAzureRMManagementLock_resourceGroupCanNotDeleteBasic(t *testing.T) {
	resourceName := "azurerm_management_lock.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMManagementLock_resourceGroupCanNotDeleteBasic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(resourceName),
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

func TestAccAzureRMManagementLock_resourceGroupCanNotDeleteComplete(t *testing.T) {
	resourceName := "azurerm_management_lock.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMManagementLock_resourceGroupCanNotDeleteComplete(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(resourceName),
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

func TestAccAzureRMManagementLock_publicIPReadOnlyBasic(t *testing.T) {
	resourceName := "azurerm_management_lock.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMManagementLock_publicIPReadOnlyBasic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(resourceName),
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

func TestAccAzureRMManagementLock_publicIPCanNotDeleteBasic(t *testing.T) {
	resourceName := "azurerm_management_lock.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccAzureRMManagementLock_publicIPCanNotDeleteBasic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(resourceName),
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

func TestAccAzureRMManagementLock_subscriptionReadOnlyBasic(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_SUBSCRIPTION_PARALLEL_LOCK")
	if !exists {
		t.Skip("`TF_ACC_SUBSCRIPTION_PARALLEL_LOCK` isn't specified - skipping since this test can't be run in Parallel")
	}

	resourceName := "azurerm_management_lock.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMManagementLock_subscriptionReadOnlyBasic(ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(resourceName),
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

func TestAccAzureRMManagementLock_subscriptionCanNotDeleteBasic(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_SUBSCRIPTION_PARALLEL_LOCK")
	if !exists {
		t.Skip("`TF_ACC_SUBSCRIPTION_PARALLEL_LOCK` isn't specified - skipping since this test can't be run in Parallel")
	}

	resourceName := "azurerm_management_lock.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMManagementLock_subscriptionCanNotDeleteBasic(ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMManagementLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementLockExists(resourceName),
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

func testCheckAzureRMManagementLockExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		scope := rs.Primary.Attributes["scope"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.LocksClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.GetByScope(ctx, scope, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Management Lock %q (Scope %q) does not exist", name, scope)
			}

			return fmt.Errorf("Bad: Get on managementLocksClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMManagementLockDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Resource.LocksClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_management_lock" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		scope := rs.Primary.Attributes["scope"]

		resp, err := client.GetByScope(ctx, scope, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}
	}

	return nil
}

func testAccAzureRMManagementLock_resourceGroupReadOnlyBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${azurerm_resource_group.test.id}"
  lock_level = "ReadOnly"
}
`, rInt, location, rInt)
}

func testAccAzureRMManagementLock_requiresImport(rInt int, location string) string {
	template := testAccAzureRMManagementLock_resourceGroupReadOnlyBasic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "import" {
  name       = "${azurerm_management_lock.test.name}"
  scope      = "${azurerm_management_lock.test.scope}"
  lock_level = "${azurerm_management_lock.test.lock_level}"
}
`, template)
}

func testAccAzureRMManagementLock_resourceGroupReadOnlyComplete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${azurerm_resource_group.test.id}"
  lock_level = "ReadOnly"
  notes      = "Hello, World!"
}
`, rInt, location, rInt)
}

func testAccAzureRMManagementLock_resourceGroupCanNotDeleteBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${azurerm_resource_group.test.id}"
  lock_level = "CanNotDelete"
}
`, rInt, location, rInt)
}

func testAccAzureRMManagementLock_resourceGroupCanNotDeleteComplete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${azurerm_resource_group.test.id}"
  lock_level = "CanNotDelete"
  notes      = "Hello, World!"
}
`, rInt, location, rInt)
}

func testAccAzureRMManagementLock_publicIPReadOnlyBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${azurerm_public_ip.test.id}"
  lock_level = "ReadOnly"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMManagementLock_publicIPCanNotDeleteBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${azurerm_public_ip.test.id}"
  lock_level = "CanNotDelete"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMManagementLock_subscriptionReadOnlyBasic(rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "current" {}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${data.azurerm_subscription.current.id}"
  lock_level = "ReadOnly"
}
`, rInt)
}

func testAccAzureRMManagementLock_subscriptionCanNotDeleteBasic(rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "current" {}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = "${data.azurerm_subscription.current.id}"
  lock_level = "CanNotDelete"
}
`, rInt)
}
