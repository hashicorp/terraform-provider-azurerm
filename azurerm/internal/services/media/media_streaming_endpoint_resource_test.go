package media_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MediaStreamingEndpointResource struct {
}

func TestAccMediaStreamingEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_endpoint", "test")
	r := MediaStreamingEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("scale_units").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaStreamingEndpoint_CDN(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_endpoint", "test")
	r := MediaStreamingEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.CDN(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("cdn_profile").HasValue("MyCDNProfile"),
				check.That(data.ResourceName).Key("cdn_provider").HasValue("StandardVerizon"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaStreamingEndpoint_MaxCacheAge(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_endpoint", "test")
	r := MediaStreamingEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.maxCacheAge(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("max_cache_age_seconds").HasValue("60"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMediaStreamingEndpoint_shouldStopWhenStarted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_endpoint", "test")
	r := MediaStreamingEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				data.CheckWithClient(r.Start),
			),
		},
	})
}

func (r MediaStreamingEndpointResource) Start(ctx context.Context, client *clients.Client, state *terraform.InstanceState) error {
	id, err := parse.StreamingEndpointID(state.ID)
	if err != nil {
		return err
	}

	future, err := client.Media.StreamingEndpointsClient.Start(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return fmt.Errorf("starting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Media.StreamingEndpointsClient.Client); err != nil {
		return fmt.Errorf("waiting for %s to start: %+v", id, err)
	}

	return nil
}

func (MediaStreamingEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.StreamingEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.StreamingEndpointsClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.StreamingEndpointProperties != nil), nil
}

func (r MediaStreamingEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_media_streaming_endpoint" "test" {
  name                        = "endpoint1"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  media_services_account_name = azurerm_media_services_account.test.name
  scale_units                 = 1
}
`, r.template(data))
}

func (r MediaStreamingEndpointResource) CDN(data acceptance.TestData) string {
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
`, r.template(data))
}

func (r MediaStreamingEndpointResource) maxCacheAge(data acceptance.TestData) string {
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
`, r.template(data))
}

func (MediaStreamingEndpointResource) template(data acceptance.TestData) string {
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
