package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SpringCloudAppRedisAssociationResource struct {
}

func TestAccSpringCloudAppRedisAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_redis_association", "test")
	r := SpringCloudAppRedisAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("redis_access_key"),
	})
}

func TestAccSpringCloudAppRedisAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_redis_association", "test")
	r := SpringCloudAppRedisAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSpringCloudAppRedisAssociation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_redis_association", "test")
	r := SpringCloudAppRedisAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("redis_access_key"),
	})
}

func TestAccSpringCloudAppRedisAssociation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_redis_association", "test")
	r := SpringCloudAppRedisAssociationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("redis_access_key"),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("redis_access_key"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("redis_access_key"),
	})
}

func (t SpringCloudAppRedisAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudAppAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppPlatform.BindingsClient.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r SpringCloudAppRedisAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_redis_association" "test" {
  name                = "acctestscarb-%d"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  redis_cache_id      = azurerm_redis_cache.test.id
  redis_access_key    = azurerm_redis_cache.test.primary_access_key
}
`, r.template(data), data.RandomInteger)
}

func (r SpringCloudAppRedisAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_redis_association" "import" {
  name                = azurerm_spring_cloud_app_redis_association.test.name
  spring_cloud_app_id = azurerm_spring_cloud_app_redis_association.test.spring_cloud_app_id
  redis_cache_id      = azurerm_spring_cloud_app_redis_association.test.redis_cache_id
  redis_access_key    = azurerm_spring_cloud_app_redis_association.test.redis_access_key
}
`, r.basic(data))
}

func (r SpringCloudAppRedisAssociationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_redis_association" "test" {
  name                = "acctestscarb-%d"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  redis_cache_id      = azurerm_redis_cache.test.id
  redis_access_key    = azurerm_redis_cache.test.secondary_access_key
  ssl_enabled         = false
}
`, r.template(data), data.RandomInteger)
}

func (r SpringCloudAppRedisAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
}

resource "azurerm_redis_cache" "test" {
  name                = "acctestredis-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  capacity            = 0
  family              = "C"
  sku_name            = "Basic"
  enable_non_ssl_port = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
