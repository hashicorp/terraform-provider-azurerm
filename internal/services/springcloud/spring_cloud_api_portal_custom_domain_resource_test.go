package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudAPIPortalCustomDomainResource struct{}

func TestAccSpringCloudAPIPortalCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_api_portal_custom_domain", "test")
	r := SpringCloudAPIPortalCustomDomainResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudAPIPortalCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_api_portal_custom_domain", "test")
	r := SpringCloudAPIPortalCustomDomainResource{}
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

func (r SpringCloudAPIPortalCustomDomainResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudAPIPortalCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.APIPortalCustomDomainClient.Get(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName, id.DomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SpringCloudAPIPortalCustomDomainResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[2]d"
  location = "%[1]s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_api_portal" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudAPIPortalCustomDomainResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_api_portal_custom_domain" "test" {
  name                       = "${azurerm_spring_cloud_service.test.name}.azuremicroservices.io"
  spring_cloud_api_portal_id = azurerm_spring_cloud_api_portal.test.id
}
`, template)
}

func (r SpringCloudAPIPortalCustomDomainResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_api_portal_custom_domain" "import" {
  name                       = azurerm_spring_cloud_api_portal_custom_domain.test.name
  spring_cloud_api_portal_id = azurerm_spring_cloud_api_portal_custom_domain.test.spring_cloud_api_portal_id
}
`, config)
}
