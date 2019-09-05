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

func TestAccAzureRMDnsSrvRecord_basic(t *testing.T) {
	resourceName := "azurerm_dns_srv_record.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMDnsSrvRecord_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists(resourceName),
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

func TestAccAzureRMDnsSrvRecord_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_dns_srv_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsSrvRecord_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMDnsSrvRecord_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_dns_srv_record"),
			},
		},
	})
}

func TestAccAzureRMDnsSrvRecord_updateRecords(t *testing.T) {
	resourceName := "azurerm_dns_srv_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsSrvRecord_basic(ri, location)
	postConfig := testAccAzureRMDnsSrvRecord_updateRecords(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsSrvRecord_withTags(t *testing.T) {
	resourceName := "azurerm_dns_srv_record.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsSrvRecord_withTags(ri, location)
	postConfig := testAccAzureRMDnsSrvRecord_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists(resourceName),
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

func testCheckAzureRMDnsSrvRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		srvName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS SRV record: %s", srvName)
		}

		conn := testAccProvider.Meta().(*ArmClient).dns.RecordSetsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, srvName, dns.SRV)
		if err != nil {
			return fmt.Errorf("Bad: Get SRV RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS SRV record %s (resource group: %s) does not exist", srvName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsSrvRecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dns.RecordSetsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_srv_record" {
			continue
		}

		srvName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, srvName, dns.SRV)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DNS SRV record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMDnsSrvRecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_srv_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }

  record {
    priority = 2
    weight   = 25
    port     = 8080
    target   = "target2.contoso.com"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsSrvRecord_requiresImport(rInt int, location string) string {
	template := testAccAzureRMDnsSrvRecord_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_srv_record" "import" {
  name                = "${azurerm_dns_srv_record.test.name}"
  resource_group_name = "${azurerm_dns_srv_record.test.resource_group_name}"
  zone_name           = "${azurerm_dns_srv_record.test.zone_name}"
  ttl                 = 300

  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }

  record {
    priority = 2
    weight   = 25
    port     = 8080
    target   = "target2.contoso.com"
  }
}
`, template)
}

func testAccAzureRMDnsSrvRecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_srv_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }

  record {
    priority = 2
    weight   = 25
    port     = 8080
    target   = "target2.contoso.com"
  }

  record {
    priority = 3
    weight   = 100
    port     = 8080
    target   = "target3.contoso.com"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsSrvRecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_srv_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }

  record {
    priority = 2
    weight   = 25
    port     = 8080
    target   = "target2.contoso.com"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsSrvRecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_srv_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }

  record {
    priority = 2
    weight   = 25
    port     = 8080
    target   = "target2.contoso.com"
  }

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}
