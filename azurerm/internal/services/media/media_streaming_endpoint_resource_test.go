package media_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
)

func TestAccAzureRMStreamingEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamingEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamingEndpoint_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "scale_units", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStreamingEndpoint_CDN(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamingEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamingEndpoint_CDN(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "cdn_profile", "MyCDNProfile"),
					resource.TestCheckResourceAttr(data.ResourceName, "cdn_provider", "StandardVerizon"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStreamingEndpoint_MaxCacheAge(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStreamingEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStreamingEndpoint_maxCacheAge(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "max_cache_age_seconds", "60"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStreamingEndpointDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Media.StreamingEndpointsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_media_streaming_endpoint" {
			continue
		}

		id, err := parse.StreamingEndpointID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Streaming Endpoint still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMStreamingEndpoint_basic(data acceptance.TestData) string {
	template := testAccAzureRMStreamingEndpoint_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_media_streaming_endpoint" "test" {
  name                        = "endpoint1"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  media_services_account_name = azurerm_media_services_account.test.name
  scale_units                 = 1
}
`, template)
}

func testAccAzureRMStreamingEndpoint_CDN(data acceptance.TestData) string {
	template := testAccAzureRMStreamingEndpoint_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_media_streaming_endpoint" "test" {
  name                        = "endpoint1"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  media_services_account_name = azurerm_media_services_account.test.name
  scale_units                 = 1
  cdn_enabled                 = true
  cdn_provider                = "StandardVerizon"
  cdn_profile                 = "MyCDNProfile"
}
`, template)
}

func testAccAzureRMStreamingEndpoint_maxCacheAge(data acceptance.TestData) string {
	template := testAccAzureRMStreamingEndpoint_template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_media_streaming_endpoint" "test" {
  name                        = "endpoint1"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  media_services_account_name = azurerm_media_services_account.test.name
  scale_units                 = 1
  access_control {
    ip_allow {
      name    = "AllowedIP"
      address = "192.168.1.1"
    }

    ip_allow {
      name    = "AnotherIp"
      address = "192.168.1.2"
    }

    akamai_signature_header_authentication_key {
      identifier = "id1"
      expiration = "2030-12-31T16:00:00Z"
      base64_key = "dGVzdGlkMQ=="
    }

    akamai_signature_header_authentication_key {
      identifier = "id2"
      expiration = "2032-01-28T16:00:00Z"
      base64_key = "dGVzdGlkMQ=="
    }
  }
  max_cache_age_seconds = 60

}
`, template)
}

func testAccAzureRMStreamingEndpoint_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-media-%d"
  location = "%s"
}
resource "azurerm_storage_account" "test" {
  name                     = "acctestsa1%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  storage_account {
    id         = azurerm_storage_account.test.id
    is_primary = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
