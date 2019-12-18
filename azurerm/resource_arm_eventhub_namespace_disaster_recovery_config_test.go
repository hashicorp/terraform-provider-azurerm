package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventHubNamespaceDisasterRecoveryConfig_basic(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace_disaster_recovery_config.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDisasterRecoveryConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceDisasterRecoveryConfig_basic(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceDisasterRecoveryConfigExists(resourceName),
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

func TestAccAzureRMEventHubNamespaceDisasterRecoveryConfig_complete(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace_disaster_recovery_config.test"
	ri := tf.AccRandTimeInt()

	// skipping due to there being no way to delete a DRC once an alternate name has been set
	// sdk bug: https://github.com/Azure/azure-sdk-for-go/issues/5893
	t.Skip()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDisasterRecoveryConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceDisasterRecoveryConfig_complete(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceDisasterRecoveryConfigExists(resourceName),
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

func TestAccAzureRMEventHubNamespaceDisasterRecoveryConfig_update(t *testing.T) {
	resourceName := "azurerm_eventhub_namespace_disaster_recovery_config.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceDisasterRecoveryConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceDisasterRecoveryConfig_basic(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceDisasterRecoveryConfigExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMEventHubNamespaceDisasterRecoveryConfig_updated(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceDisasterRecoveryConfigExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMEventHubNamespaceDisasterRecoveryConfig_updated_removed(ri, testLocation(), testAltLocation()),
			},
		},
	})
}

func testCheckAzureRMEventHubNamespaceDisasterRecoveryConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Eventhub.DisasterRecoveryConfigsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_namespace_disaster_recovery_config" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: EventHub Namespace Disaster Recovery Configs %q (namespace %q / resource group: %q): %+v", name, namespaceName, resourceGroup, err)
			}
		}
	}

	return nil
}

func testCheckAzureRMEventHubNamespaceDisasterRecoveryConfigExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).Eventhub.DisasterRecoveryConfigsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: EventHub Namespace Disaster Recovery Configs %q (namespace %q / resource group: %q) does not exist", name, namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: EventHub Namespace Disaster Recovery Configs %q (namespace %q / resource group: %q): %+v", name, namespaceName, resourceGroup, err)
		}

		return nil
	}
}

func testAccAzureRMEventHubNamespaceDisasterRecoveryConfig_basic(rInt int, location string, altlocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "testa" {
  name                = "acctest-EHN-%[1]d-a"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testb" {
  name                = "acctest-EHN-%[1]d-b"
  location            = "%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "test" {
  name                 = "acctest-EHN-DRC-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  namespace_name       = "${azurerm_eventhub_namespace.testa.name}"
  partner_namespace_id = "${azurerm_eventhub_namespace.testb.id}"
}

`, rInt, location, altlocation)
}

func testAccAzureRMEventHubNamespaceDisasterRecoveryConfig_complete(rInt int, location string, altlocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "testa" {
  name                = "acctest-EHN-%[1]d-a"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testb" {
  name                = "acctest-EHN-%[1]d-b"
  location            = "%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "test" {
  name                 = "${azurerm_eventhub_namespace.testa.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  namespace_name       = "${azurerm_eventhub_namespace.testa.name}"
  partner_namespace_id = "${azurerm_eventhub_namespace.testb.id}"
  alternate_name       = "acctest-EHN-DRC-%[1]d-alt"
}

`, rInt, location, altlocation)
}

func testAccAzureRMEventHubNamespaceDisasterRecoveryConfig_updated(rInt int, location string, altlocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "testa" {
  name                = "acctest-EHN-%[1]d-a"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testb" {
  name                = "acctest-EHN-%[1]d-b"
  location            = "%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testc" {
  name                = "acctest-EHN-%[1]d-c"
  location            = "%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace_disaster_recovery_config" "test" {
  name                 = "acctest-EHN-DRC-%[1]d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  namespace_name       = "${azurerm_eventhub_namespace.testa.name}"
  partner_namespace_id = "${azurerm_eventhub_namespace.testc.id}"
}

`, rInt, location, altlocation)
}

func testAccAzureRMEventHubNamespaceDisasterRecoveryConfig_updated_removed(rInt int, location string, altlocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "testa" {
  name                = "acctest-EHN-%[1]d-a"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testb" {
  name                = "acctest-EHN-%[1]d-b"
  location            = "%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace" "testc" {
  name                = "acctest-EHN-%[1]d-c"
  location            = "%[3]s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}
`, rInt, location, altlocation)
}
