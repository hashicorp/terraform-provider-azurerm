package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventGridDomain_basic(t *testing.T) {
	resourceName := "azurerm_eventgrid_domain.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridDomain_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridDomainExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_access_key"),
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

func TestAccAzureRMEventGridDomain_mapping(t *testing.T) {
	resourceName := "azurerm_eventgrid_domain.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridDomain_mapping(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "input_mapping_fields.0.topic", "test"),
					resource.TestCheckResourceAttr(resourceName, "input_mapping_fields.0.topic", "test"),
					resource.TestCheckResourceAttr(resourceName, "input_mapping_default_values.0.data_version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "input_mapping_default_values.0.subject", "DefaultSubject"),
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

func TestAccAzureRMEventGridDomain_basicWithTags(t *testing.T) {
	resourceName := "azurerm_eventgrid_domain.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventGridDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventGridDomain_basicWithTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventGridDomainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
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

func testCheckAzureRMEventGridDomainDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.DomainsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventgrid_domain" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("EventGrid Domain still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMEventGridDomainExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for EventGrid Domain: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).EventGrid.DomainsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: EventGrid Domain %q (resource group: %s) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on eventGridDomainsClient: %s", err)
		}

		return nil
	}
}

func testAccAzureRMEventGridDomain_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMEventGridDomain_mapping(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  input_schema = "CustomEventSchema"

  input_mapping_fields {
    topic      = "test"
    event_type = "test"
  }

  input_mapping_default_values {
    data_version = "1.0"
    subject      = "DefaultSubject"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMEventGridDomain_basicWithTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    "foo" = "bar"
  }
}
`, rInt, location, rInt)
}
