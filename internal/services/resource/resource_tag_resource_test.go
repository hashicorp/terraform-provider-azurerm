package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceTagResource struct{}

func TestAccResourceTag_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_tag", "test")
	testResource := ResourceTagResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.ImportStep(),
	})
}

func TestAccResourceTag_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_tag", "test")
	testResource := ResourceTagResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       testResource.basicConfig,
			TestResource: testResource,
		}),
	})
}

func TestAccResourceTag_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_tag", "test")
	testResource := ResourceTagResource{}
	assert := check.That(data.ResourceName)
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		{
			Config: testResource.withTagsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("key").HasValue("owner"),
				assert.Key("value").HasValue("Terraform"),
			),
		},
		data.ImportStep(),
		{
			Config: testResource.withTagsUpdatedConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				assert.ExistsInAzure(testResource),
				assert.Key("key").HasValue("owner"),
				assert.Key("value").HasValue("Human"),
			),
		},
		data.ImportStep(),
	})
}

func (t ResourceTagResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	resourceId := state.Attributes["resource_id"]
	key := state.Attributes["key"]

	resp, err := client.Resource.TagsClient.GetAtScope(ctx, resourceId)
	if err != nil {
		return nil, fmt.Errorf("retrieving tags at scope %q: %+v", resourceId, err)
	}

	delete(resp.Properties.Tags, key)

	_, err = client.Resource.TagsClient.CreateOrUpdateAtScope(ctx, resourceId, resp)

	if err != nil {
		return nil, fmt.Errorf("remove resourceTag %q at scope %q: %+v", key, resourceId, err)
	}

	return utils.Bool(true), nil
}

func (t ResourceTagResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	resourceId := state.Attributes["resource_id"]
	key := state.Attributes["key"]

	resp, err := client.Resource.TagsClient.GetAtScope(ctx, resourceId)
	if err != nil {
		return nil, fmt.Errorf("retrieving tags at scope %q: %+v", resourceId, err)
	}

	_, ok := resp.Properties.Tags[key]
	return utils.Bool(ok), nil
}

func (t ResourceTagResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  lifecycle {
    ignore_changes = [tags["owner"]]
  }
}

resource "azurerm_resource_tag" "test" {
  resource_id = azurerm_resource_group.test.id

  key   = "owner"
  value = "Terraform"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (t ResourceTagResource) withTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }

  lifecycle {
    ignore_changes = [tags["owner"]]
  }
}

resource "azurerm_resource_tag" "test" {
  resource_id = azurerm_resource_group.test.id

  key   = "owner"
  value = "Terraform"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (t ResourceTagResource) withTagsUpdatedConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }

  lifecycle {
    ignore_changes = [tags["owner"]]
  }
}

resource "azurerm_resource_tag" "test" {
  resource_id = azurerm_resource_group.test.id

  key   = "owner"
  value = "Human"
}
`, data.RandomInteger, data.Locations.Primary)
}
