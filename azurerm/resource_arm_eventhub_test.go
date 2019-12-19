package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMEventHubPartitionCount_validation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{
			Value:    0,
			ErrCount: 1,
		},
		{
			Value:    1,
			ErrCount: 0,
		},
		{
			Value:    2,
			ErrCount: 0,
		},
		{
			Value:    3,
			ErrCount: 0,
		},
		{
			Value:    21,
			ErrCount: 0,
		},
		{
			Value:    32,
			ErrCount: 0,
		},
		{
			Value:    33,
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateEventHubPartitionCount(tc.Value, "azurerm_eventhub")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM EventHub Partition Count to trigger a validation error")
		}
	}
}

func TestAccAzureRMEventHubMessageRetentionCount_validation(t *testing.T) {
	cases := []struct {
		Value    int
		ErrCount int
	}{
		{
			Value:    0,
			ErrCount: 1,
		}, {
			Value:    1,
			ErrCount: 0,
		}, {
			Value:    2,
			ErrCount: 0,
		}, {
			Value:    3,
			ErrCount: 0,
		}, {
			Value:    4,
			ErrCount: 0,
		}, {
			Value:    5,
			ErrCount: 0,
		}, {
			Value:    6,
			ErrCount: 0,
		}, {
			Value:    7,
			ErrCount: 0,
		}, {
			Value:    8,
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateEventHubMessageRetentionCount(tc.Value, "azurerm_eventhub")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM EventHub Message Retention Count to trigger a validation error")
		}
	}
}

func TestAccAzureRMEventHubArchiveNameFormat_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 9,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 0,
		},
		{
			Value:    "Prod_{Eventub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Month}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}/{Day}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Hour}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Minute}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Second}",
			ErrCount: 1,
		},
		{
			Value:    "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateEventHubArchiveNameFormat(tc.Value, "azurerm_eventhub")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to trigger a validation error", tc.Value)
		}
	}
}

func TestAccAzureRMEventHub_basic(t *testing.T) {
	resourceName := "azurerm_eventhub.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHub_basic(ri, acceptance.Location(), 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
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

func TestAccAzureRMEventHub_basicOnePartition(t *testing.T) {
	resourceName := "azurerm_eventhub.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHub_basic(ri, acceptance.Location(), 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "partition_count", "1"),
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

func TestAccAzureRMEventHub_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_eventhub.test"
	ri := tf.AccRandTimeInt()

	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHub_basic(ri, location, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMEventHub_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_eventhub"),
			},
		},
	})
}

func TestAccAzureRMEventHub_partitionCountUpdate(t *testing.T) {
	resourceName := "azurerm_eventhub.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMEventHub_basic(ri, acceptance.Location(), 2)
	postConfig := testAccAzureRMEventHub_partitionCountUpdate(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "partition_count", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "partition_count", "10"),
				),
			},
		},
	})
}

func TestAccAzureRMEventHub_standard(t *testing.T) {
	resourceName := "azurerm_eventhub.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHub_standard(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
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

func TestAccAzureRMEventHub_captureDescription(t *testing.T) {
	resourceName := "azurerm_eventhub.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHub_captureDescription(ri, rs, acceptance.Location(), true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "capture_description.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "capture_description.0.skip_empty_archives", "true"),
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

func TestAccAzureRMEventHub_captureDescriptionDisabled(t *testing.T) {
	resourceName := "azurerm_eventhub.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(5)
	location := acceptance.Location()

	config := testAccAzureRMEventHub_captureDescription(ri, rs, location, true)
	updatedConfig := testAccAzureRMEventHub_captureDescription(ri, rs, location, false)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "capture_description.0.enabled", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "capture_description.0.enabled", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMEventHub_messageRetentionUpdate(t *testing.T) {
	resourceName := "azurerm_eventhub.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMEventHub_standard(ri, acceptance.Location())
	postConfig := testAccAzureRMEventHub_messageRetentionUpdate(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "message_retention", "7"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "message_retention", "5"),
				),
			},
		},
	})
}

func testCheckAzureRMEventHubDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.EventHubsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, namespaceName, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("EventHub still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

func testCheckAzureRMEventHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Event Hub: %s", name)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.EventHubsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := conn.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on eventHubClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Event Hub %q (namespace %q / resource group: %q) does not exist", name, namespaceName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMEventHub_basic(rInt int, location string, partitionCount int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = %d
  message_retention   = 1
}
`, rInt, location, rInt, rInt, partitionCount)
}

func testAccAzureRMEventHub_requiresImport(rInt int, location string) string {
	template := testAccAzureRMEventHub_basic(rInt, location, 2)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub" "import" {
  name                = "${azurerm_eventhub.test.name}"
  namespace_name      = "${azurerm_eventhub.test.namespace_name}"
  resource_group_name = "${azurerm_eventhub.test.resource_group_name}"
  partition_count     = "${azurerm_eventhub.test.partition_count}"
  message_retention   = "${azurerm_eventhub.test.message_retention}"
}
`, template)
}

func testAccAzureRMEventHub_partitionCountUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 10
  message_retention   = 1
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMEventHub_standard(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 7
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMEventHub_captureDescription(rInt int, rString string, location string, enabled bool) string {
	enabledString := strconv.FormatBool(enabled)
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 7

  capture_description {
    enabled             = %s
    encoding            = "Avro"
    interval_in_seconds = 60
    size_limit_in_bytes = 10485760
    skip_empty_archives = true

    destination {
      name                = "EventHubArchive.AzureBlockBlob"
      archive_name_format = "Prod_{EventHub}/{Namespace}\\{PartitionId}_{Year}_{Month}/{Day}/{Hour}/{Minute}/{Second}"
      blob_container_name = "${azurerm_storage_container.test.name}"
      storage_account_id  = "${azurerm_storage_account.test.id}"
    }
  }
}
`, rInt, location, rString, rInt, rInt, enabledString)
}

func testAccAzureRMEventHub_messageRetentionUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-eventhub-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 5
}
`, rInt, location, rInt, rInt)
}
