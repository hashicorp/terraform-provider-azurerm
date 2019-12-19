package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMPrivateDnsSrvRecord_basic(t *testing.T) {
	resourceName := "azurerm_private_dns_srv_record.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPrivateDnsSrvRecord_basic(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsSrvRecordExists(resourceName),
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

func TestAccAzureRMPrivateDnsSrvRecord_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_private_dns_srv_record.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsSrvRecord_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsSrvRecordExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMPrivateDnsSrvRecord_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_private_dns_srv_record"),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsSrvRecord_updateRecords(t *testing.T) {
	resourceName := "azurerm_private_dns_srv_record.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMPrivateDnsSrvRecord_basic(ri, location)
	postConfig := testAccAzureRMPrivateDnsSrvRecord_updateRecords(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsSrvRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsSrvRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsSrvRecord_withTags(t *testing.T) {
	resourceName := "azurerm_private_dns_srv_record.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	preConfig := testAccAzureRMPrivateDnsSrvRecord_withTags(ri, location)
	postConfig := testAccAzureRMPrivateDnsSrvRecord_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsSrvRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsSrvRecordExists(resourceName),
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

func testCheckAzureRMPrivateDnsSrvRecordExists(resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: no resource group found in state for Private DNS SRV record: %s", srvName)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.SRV, srvName)
		if err != nil {
			return fmt.Errorf("Bad: Get SRV RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Private DNS SRV record %s (resource group: %s) does not exist", srvName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPrivateDnsSrvRecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_dns_srv_record" {
			continue
		}

		srvName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.SRV, srvName)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Private DNS SRV record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMPrivateDnsSrvRecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-prvdns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_srv_record" "test" {
  name                = "testaccsrv%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }

  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsSrvRecord_requiresImport(rInt int, location string) string {
	template := testAccAzureRMPrivateDnsSrvRecord_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_srv_record" "import" {
  name                = "${azurerm_private_dns_srv_record.test.name}"
  resource_group_name = "${azurerm_private_dns_srv_record.test.resource_group_name}"
  zone_name           = "${azurerm_private_dns_srv_record.test.zone_name}"
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }
  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }
}
`, template)
}

func testAccAzureRMPrivateDnsSrvRecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}
	
resource "azurerm_private_dns_zone" "test" {
	name                = "testzone%d.com"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_srv_record" "test" {
  name                = "test%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }
  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }
  record {
    priority = 20
    weight   = 100
    port     = 8080
    target   = "target3.contoso.com"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMPrivateDnsSrvRecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}
  
resource "azurerm_private_dns_zone" "test" {
	name                = "testzone%d.com"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_srv_record" "test" {
  name                = "test%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }
  record {
    priority = 10
    weight   = 10
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

func testAccAzureRMPrivateDnsSrvRecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
	name                = "testzone%d.com"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_srv_record" "test" {
  name                = "test%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_private_dns_zone.test.name}"
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }
  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}
