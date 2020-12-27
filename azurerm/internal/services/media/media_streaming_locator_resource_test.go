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

type StreamingLocatorResource struct {}

func TestAccStreamingLocator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_locator", "test")
	r := StreamingLocatorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Locator-1"),
				check.That(data.ResourceName).Key("asset_name").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamingLocator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_locator", "test")
	r := StreamingLocatorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Locator-1"),
				check.That(data.ResourceName).Key("asset_name").HasValue("test"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStreamingLocator_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_locator", "test")
	r := StreamingLocatorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("start_time").HasValue("2018-03-01T00:00:00Z"),
				check.That(data.ResourceName).Key("end_time").HasValue("2028-12-31T23:59:59Z"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamingLocator_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_streaming_locator", "test")
	r := StreamingLocatorResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Locator-1"),
				check.That(data.ResourceName).Key("asset_name").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("start_time").HasValue("2018-03-01T00:00:00Z"),
				check.That(data.ResourceName).Key("end_time").HasValue("2028-12-31T23:59:59Z"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("Locator-1"),
				check.That(data.ResourceName).Key("asset_name").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func (StreamingLocatorResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.StreamingLocatorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Media.StreamingLocatorsClient.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Content Key Policy %s (Media Services Account %s) (resource group: %s): %v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.StreamingLocatorProperties != nil), nil
}

func (r StreamingLocatorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_streaming_locator" "test" {
  name                        = "Locator-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  streaming_policy_name       = "Predefined_ClearStreamingOnly"
  asset_name                  = azurerm_media_asset.test.name
}

`, r.template(data))
}

func (r StreamingLocatorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_streaming_locator" "import" {
  name                        = azurerm_media_streaming_locator.test.name
  resource_group_name         = azurerm_media_streaming_locator.test.resource_group_name
  media_services_account_name = azurerm_media_streaming_locator.test.media_services_account_name
  streaming_policy_name       = "Predefined_ClearStreamingOnly"
  asset_name                  = azurerm_media_asset.test.name
}

`, r.basic(data))
}

func (r StreamingLocatorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_media_streaming_locator" "test" {
  name                        = "Job-1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  streaming_policy_name       = "Predefined_DownloadOnly"
  asset_name                  = azurerm_media_asset.test.name
  start_time                  = "2018-03-01T00:00:00Z"
  end_time                    = "2028-12-31T23:59:59Z"
  streaming_locator_id        = "90000000-0000-0000-0000-000000000000"
  alternative_media_id        = "my-Alternate-MediaID"
}

`, r.template(data))
}

func (StreamingLocatorResource) template(data acceptance.TestData) string {
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

resource "azurerm_media_asset" "test" {
  name                        = "test"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
