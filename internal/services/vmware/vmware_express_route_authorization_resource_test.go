package vmware_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2020-03-20/authorizations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VmwareExpressRouteAuthorizationResource struct{}

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
	id, err := authorizations.ParseAuthorizationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Vmware.AuthorizationClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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
