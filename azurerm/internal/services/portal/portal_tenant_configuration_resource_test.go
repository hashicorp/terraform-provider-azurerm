package portal_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/portal/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PortalTenantConfigurationResource struct{}

func TestAccPortalTenantConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_tenant_configuration", "test")
	r := PortalTenantConfigurationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPortalTenantConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_tenant_configuration", "test")
	r := PortalTenantConfigurationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r PortalTenantConfigurationResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.TenantConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Portal.TenantConfigurationsClient.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.ConfigurationProperties != nil), nil
}

func (r PortalTenantConfigurationResource) basic() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_tenant_configuration" "test" {
  enforce_private_markdown_storage = true
}
`)
}

func (r PortalTenantConfigurationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_portal_tenant_configuration" "import" {
  enforce_private_markdown_storage = azurerm_portal_tenant_configuration.test.enforce_private_markdown_storage
}
`, r.basic())
}
