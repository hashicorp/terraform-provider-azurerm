package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDnsPtrRecord_basic(t *testing.T) {
	resourceName := "azurerm_dns_ptr_record.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMDnsPtrRecord_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsPtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists(resourceName),
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

func TestAccAzureRMDnsPtrRecord_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_dns_ptr_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsPtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsPtrRecord_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMDnsPtrRecord_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_dns_ptr_record"),
			},
		},
	})
}

func TestAccAzureRMDnsPtrRecord_updateRecords(t *testing.T) {
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsPtrRecord_basic(ri, location)
	postConfig := testAccAzureRMDnsPtrRecord_updateRecords(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsPtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists("azurerm_dns_ptr_record.test"),
					resource.TestCheckResourceAttr("azurerm_dns_ptr_record.test", "records.#", "2"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists("azurerm_dns_ptr_record.test"),
					resource.TestCheckResourceAttr("azurerm_dns_ptr_record.test", "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsPtrRecord_withTags(t *testing.T) {
	resourceName := "azurerm_dns_ptr_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsPtrRecord_withTags(ri, location)
	postConfig := testAccAzureRMDnsPtrRecord_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsPtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
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

func testCheckAzureRMDnsPtrRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		ptrName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS PTR record: %s", ptrName)
		}

		conn := testAccProvider.Meta().(*ArmClient).dns.RecordSetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, ptrName, dns.PTR)
		if err != nil {
			return fmt.Errorf("Bad: Get PTR RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS PTR record %s (resource group: %s) does not exist", ptrName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsPtrRecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dns.RecordSetsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_ptr_record" {
			continue
		}

		ptrName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, ptrName, dns.PTR)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DNS PTR record still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDnsPtrRecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ptr_record" "test" {
  name                = "testptrrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["hashicorp.com", "microsoft.com"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsPtrRecord_requiresImport(rInt int, location string) string {
	template := testAccAzureRMDnsPtrRecord_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_ptr_record" "import" {
  name                = "${azurerm_dns_ptr_record.test.name}"
  resource_group_name = "${azurerm_dns_ptr_record.test.resource_group_name}"
  zone_name           = "${azurerm_dns_ptr_record.test.zone_name}"
  ttl                 = 300
  records             = ["hashicorp.com", "microsoft.com"]
}
`, template)
}

func testAccAzureRMDnsPtrRecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ptr_record" "test" {
  name                = "testptrrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["hashicorp.com", "microsoft.com", "reddit.com"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsPtrRecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ptr_record" "test" {
  name                = "testptrrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["hashicorp.com", "microsoft.com"]

  tags = {
    environment = "Dev"
    cost_center = "Ops"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsPtrRecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ptr_record" "test" {
  name                = "testptrrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["hashicorp.com", "microsoft.com"]

  tags = {
    environment = "Stage"
  }
}
`, rInt, location, rInt, rInt)
}
