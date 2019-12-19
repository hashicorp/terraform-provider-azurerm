package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDnsARecord_basic(t *testing.T) {
	resourceName := "azurerm_dns_a_record.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMDnsARecord_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "fqdn"),
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

func TestAccAzureRMDnsARecord_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_dns_a_record.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsARecord_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMDnsARecord_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_dns_a_record"),
			},
		},
	})
}

func TestAccAzureRMDnsARecord_updateRecords(t *testing.T) {
	resourceName := "azurerm_dns_a_record.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMDnsARecord_basic(ri, location)
	postConfig := testAccAzureRMDnsARecord_updateRecords(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsARecord_withTags(t *testing.T) {
	resourceName := "azurerm_dns_a_record.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMDnsARecord_withTags(ri, location)
	postConfig := testAccAzureRMDnsARecord_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
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

func TestAccAzureRMDnsARecord_withAlias(t *testing.T) {
	resourceName := "azurerm_dns_a_record.test"
	targetResourceName := "azurerm_public_ip.test"
	targetResourceName2 := "azurerm_public_ip.test2"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMDnsARecord_withAlias(ri, location)
	postConfig := testAccAzureRMDnsARecord_withAliasUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "target_resource_id", targetResourceName, "id"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "target_resource_id", targetResourceName2, "id"),
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

func TestAccAzureRMDnsARecord_RecordsToAlias(t *testing.T) {
	resourceName := "azurerm_dns_a_record.test"
	targetResourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMDnsARecord_AliasToRecordsUpdate(ri, location)
	postConfig := testAccAzureRMDnsARecord_AliasToRecords(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "target_resource_id", targetResourceName, "id"),
					resource.TestCheckNoResourceAttr(resourceName, "records"),
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

func TestAccAzureRMDnsARecord_AliasToRecords(t *testing.T) {
	resourceName := "azurerm_dns_a_record.test"
	targetResourceName := "azurerm_public_ip.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMDnsARecord_AliasToRecords(ri, location)
	postConfig := testAccAzureRMDnsARecord_AliasToRecordsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "target_resource_id", targetResourceName, "id"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "target_resource_id"),
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

func testCheckAzureRMDnsARecordExists(resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for DNS A record: %s", aName)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, aName, dns.A)
		if err != nil {
			return fmt.Errorf("Bad: Get A RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS A record %s (resource group: %s) does not exist", aName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsARecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_a_record" {
			continue
		}

		aName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, aName, dns.A)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DNS A record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMDnsARecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsARecord_requiresImport(rInt int, location string) string {
	template := testAccAzureRMDnsARecord_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_a_record" "import" {
  name                = "${azurerm_dns_a_record.test.name}"
  resource_group_name = "${azurerm_dns_a_record.test.resource_group_name}"
  zone_name           = "${azurerm_dns_a_record.test.zone_name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]
}
`, template)
}

func testAccAzureRMDnsARecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5", "1.2.3.7"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsARecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsARecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsARecord_withAlias(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_public_ip" "test" {
	name                = "mypublicip%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	allocation_method   = "Dynamic"
	ip_version          = "IPv4"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  target_resource_id   = "${azurerm_public_ip.test.id}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMDnsARecord_withAliasUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_public_ip" "test2" {
	name                = "mypublicip%d2"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	allocation_method   = "Dynamic"
	ip_version          = "IPv4"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  target_resource_id  = "${azurerm_public_ip.test2.id}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMDnsARecord_AliasToRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_public_ip" "test" {
	name                = "mypublicip%d"
	location            = "${azurerm_resource_group.test.location}"
	resource_group_name = "${azurerm_resource_group.test.name}"
	allocation_method   = "Dynamic"
	ip_version          = "IPv4"
}

resource "azurerm_dns_a_record" "test" {
	name                = "myarecord%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
	zone_name           = "${azurerm_dns_zone.test.name}"
	ttl                 = 300
	target_resource_id  = "${azurerm_public_ip.test.id}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMDnsARecord_AliasToRecordsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]
}
`, rInt, location, rInt, rInt)
}
