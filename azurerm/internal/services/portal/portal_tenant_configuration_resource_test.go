package portal_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/portal/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PortalTenantConfigurationResource struct{}

func TestAccPortalTenantConfiguration(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being able provision one default Tenant Configuration at a time
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"resource": {
			"basic":          testAccPortalTenantConfiguration_basic,
			"update":         testAccPortalTenantConfiguration_update,
			"requiresImport": testAccPortalTenantConfiguration_requiresImport,
		},
	})
}

func testAccPortalTenantConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_portal_tenant_configuration", "test")
	r := PortalTenantConfigurationResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccPortalTenantConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_portal_tenant_configuration", "test")
	r := PortalTenantConfigurationResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccPortalTenantConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_portal_tenant_configuration", "test")
	r := PortalTenantConfigurationResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PortalTenantConfigurationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PortalTenantConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Portal.TenantConfigurationsClient.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.ConfigurationProperties != nil), nil
}

func (r PortalTenantConfigurationResource) basic(enforcePrivateMarkdownStorage bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_portal_tenant_configuration" "test" {
  private_markdown_storage_enforced = %t
}
`, enforcePrivateMarkdownStorage)
}

func (r PortalTenantConfigurationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_portal_tenant_configuration" "import" {
  private_markdown_storage_enforced = azurerm_portal_tenant_configuration.test.private_markdown_storage_enforced
}
`, r.basic(true))
}
