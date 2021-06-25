package vmware_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VmwareExpressRouteAuthorizationResource struct {
}

func TestAccVmwareExpressRouteAuthorization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vmware_express_route_authorization", "test")
	r := VmwareExpressRouteAuthorizationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("express_route_authorization_id").Exists(),
				check.That(data.ResourceName).Key("express_route_authorization_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVmwareExpressRouteAuthorization_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vmware_express_route_authorization", "test")
	r := VmwareExpressRouteAuthorizationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (VmwareExpressRouteAuthorizationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ExpressRouteAuthorizationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Vmware.AuthorizationClient.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.AuthorizationName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return utils.Bool(resp.ExpressRouteAuthorizationProperties != nil), nil
}

func (r VmwareExpressRouteAuthorizationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vmware_express_route_authorization" "test" {
  name             = "acctest-VmwareAuthorization-%d"
  private_cloud_id = azurerm_vmware_private_cloud.test.id
}
`, VmwarePrivateCloudResource{}.basic(data), data.RandomInteger)
}

func (r VmwareExpressRouteAuthorizationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vmware_express_route_authorization" "import" {
  name             = azurerm_vmware_express_route_authorization.test.name
  private_cloud_id = azurerm_vmware_private_cloud.test.id
}
`, r.basic(data))
}
