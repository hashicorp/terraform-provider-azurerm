package vmware_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VmwareAuthorizationResource struct {
}

func TestAccVmwareAuthorization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vmware_authorization", "test")
	r := VmwareAuthorizationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("express_route_authorization_id").Exists(),
				check.That(data.ResourceName).Key("express_route_authorization_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVmwareAuthorization_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vmware_authorization", "test")
	r := VmwareAuthorizationResource{}

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

func (VmwareAuthorizationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AuthorizationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Vmware.AuthorizationClient.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return utils.Bool(resp.ExpressRouteAuthorizationProperties != nil), nil
}

func (r VmwareAuthorizationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vmware_authorization" "test" {
  name             = "acctest-VmwareAuthorization-%d"
  private_cloud_id = azurerm_avs_private_cloud.test.id
}
`, VmwarePrivateCloudResource{}.basic(data), data.RandomInteger)
}

func (r VmwareAuthorizationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vmware_authorization" "import" {
  name             = azurerm_vmware_authorization.test.name
  private_cloud_id = azurerm_avs_private_cloud.test.id
}
`, r.basic(data))
}
