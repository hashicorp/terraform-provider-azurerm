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

func TestAccAzureRMPrivateDnsCNameRecord_basic(t *testing.T) {
	resourceName := "azurerm_private_dns_cname_record.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPrivateDnsCNameRecord_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(resourceName),
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

func TestAccAzureRMPrivateDnsCNameRecord_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_private_dns_cname_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsCNameRecord_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMPrivateDnsCNameRecord_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_private_dns_cname_record"),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsCNameRecord_subdomain(t *testing.T) {
	resourceName := "azurerm_private_dns_cname_record.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPrivateDnsCNameRecord_subdomain(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record", "test.contoso.com"),
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

func TestAccAzureRMPrivateDnsCNameRecord_updateRecords(t *testing.T) {
	resourceName := "azurerm_private_dns_cname_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMPrivateDnsCNameRecord_basic(ri, location)
	postConfig := testAccAzureRMPrivateDnsCNameRecord_updateRecords(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(resourceName),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsCNameRecord_withTags(t *testing.T) {
	resourceName := "azurerm_private_dns_cname_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMPrivateDnsCNameRecord_withTags(ri, location)
	postConfig := testAccAzureRMPrivateDnsCNameRecord_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(resourceName),
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

func testCheckAzureRMPrivateDnsCNameRecordExists(resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for Private DNS CNAME record: %s", aName)
		}

		conn := testAccProvider.Meta().(*ArmClient).privateDns.RecordSetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.CNAME, aName)
		if err != nil {
			return fmt.Errorf("Bad: Get CNAME RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Private DNS CNAME record %s (resource group: %s) does not exist", aName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPrivateDnsCNameRecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).privateDns.RecordSetsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_dns_cname_record" {
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

		return fmt.Errorf("Private DNS CNAME record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMPrivateDnsCNameRecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record              = "contoso.com"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsCNameRecord_requiresImport(rInt int, location string) string {
	template := testAccAzureRMPrivateDnsCNameRecord_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_cname_record" "import" {
  name                = "${azurerm_private_dns_cname_record.test.name}"
  resource_group_name = "${azurerm_private_dns_cname_record.test.resource_group_name}"
  zone_name           = "${azurerm_private_dns_cname_record.test.zone_name}"
  ttl                 = 300
  record              = "contoso.com"
}
`, template)
}

func testAccAzureRMPrivateDnsCNameRecord_subdomain(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record              = "test.contoso.com"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsCNameRecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record              = "contoso.com"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsCNameRecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record              = "contoso.com"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsCNameRecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record              = "contoso.com"

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}
