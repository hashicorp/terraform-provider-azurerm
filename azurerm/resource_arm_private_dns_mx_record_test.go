package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMPrivateDnsMxRecord_basic(t *testing.T) {
	resourceName := "azurerm_private_dns_mx_record.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPrivateDnsMxRecord_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(resourceName),
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

func TestAccAzureRMPrivateDnsMxRecord_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_private_dns_mx_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsMxRecord_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMPrivateDnsMxRecord_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_private_dns_mx_record"),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsMxRecord_updateRecords(t *testing.T) {
	resourceName := "azurerm_private_dns_mx_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMPrivateDnsMxRecord_basic(ri, location)
	postConfig := testAccAzureRMPrivateDnsMxRecord_updateRecords(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsMxRecord_withTags(t *testing.T) {
	resourceName := "azurerm_private_dns_mx_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMPrivateDnsMxRecord_withTags(ri, location)
	postConfig := testAccAzureRMPrivateDnsMxRecord_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(resourceName),
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

func TestAccAzureRMPrivateDnsMxRecord_emptyName(t *testing.T) {
	resourceName := "azurerm_private_dns_mx_record.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPrivateDnsMxRecord_emptyName(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(resourceName),
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

func testCheckAzureRMPrivateDnsMxRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		mxName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Private DNS MX record: %s", mxName)
		}

		conn := testAccProvider.Meta().(*ArmClient).PrivateDns.RecordSetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.MX, mxName)
		if err != nil {
			return fmt.Errorf("Bad: Get MX RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Private DNS MX record %s (resource group: %s) does not exist", mxName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPrivateDnsMxRecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).PrivateDns.RecordSetsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_dns_mx_record" {
			continue
		}

		mxName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.MX, mxName)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Private DNS MX record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMPrivateDnsMxRecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-prvdns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_mx_record" "test" {
  name                = "testaccmx%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }

  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsMxRecord_emptyName(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-prvdns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_mx_record" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }

  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMPrivateDnsMxRecord_requiresImport(rInt int, location string) string {
	template := testAccAzureRMPrivateDnsMxRecord_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_mx_record" "import" {
  name                = "${azurerm_private_dns_mx_record.test.name}"
  resource_group_name = "${azurerm_private_dns_mx_record.test.resource_group_name}"
  zone_name           = "${azurerm_private_dns_mx_record.test.zone_name}"
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }
  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }
}
`, template)
}

func testAccAzureRMPrivateDnsMxRecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-prvdns-%d"
	location = "%s"
}
	
resource "azurerm_private_dns_zone" "test" {
	name                = "testzone%d.com"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_mx_record" "test" {
  name                = "testaccmx%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }
  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }
  record {
    preference = 20
    exchange   = "backupmx.contoso.com"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsMxRecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-prvdns-%d"
	location = "%s"
}
  
resource "azurerm_private_dns_zone" "test" {
	name                = "testzone%d.com"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_mx_record" "test" {
  name                = "testaccmx%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }
  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsMxRecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-prvdns-%d"
	location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
	name                = "testzone%d.com"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_mx_record" "test" {
  name                = "testaccmx%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }
  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}
