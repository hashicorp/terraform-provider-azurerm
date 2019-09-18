package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMPrivateDnsARecord_basic(t *testing.T) {
	resourceName := "azurerm_private_dns_a_record.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPrivateDnsARecord_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsARecordExists(resourceName),
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

func TestAccAzureRMPrivateDnsARecord_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_private_dns_a_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsARecord_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsARecordExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMPrivateDnsARecord_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_private_dns_a_record"),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsARecord_updateRecords(t *testing.T) {
	resourceName := "azurerm_private_dns_a_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMPrivateDnsARecord_basic(ri, location)
	postConfig := testAccAzureRMPrivateDnsARecord_updateRecords(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsARecord_withTags(t *testing.T) {
	resourceName := "azurerm_private_dns_a_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMPrivateDnsARecord_withTags(ri, location)
	postConfig := testAccAzureRMPrivateDnsARecord_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsARecordExists(resourceName),
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

func testCheckAzureRMPrivateDnsARecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		aName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Private DNS A record: %s", aName)
		}

		conn := testAccProvider.Meta().(*ArmClient).privateDns.RecordSetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.A, aName)
		if err != nil {
			return fmt.Errorf("Bad: Get A RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Private DNS A record %s (resource group: %s) does not exist", aName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPrivateDnsARecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).privateDns.RecordSetsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_dns_a_record" {
			continue
		}

		aName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.A, aName)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Private DNS A record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMPrivateDnsARecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsARecord_requiresImport(rInt int, location string) string {
	template := testAccAzureRMPrivateDnsARecord_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_a_record" "import" {
  name                = "${azurerm_private_dns_a_record.test.name}"
  resource_group_name = "${azurerm_private_dns_a_record.test.resource_group_name}"
  zone_name           = "${azurerm_private_dns_a_record.test.zone_name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]
}
`, template)
}

func testAccAzureRMPrivateDnsARecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5", "1.2.3.7"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsARecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsARecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}
