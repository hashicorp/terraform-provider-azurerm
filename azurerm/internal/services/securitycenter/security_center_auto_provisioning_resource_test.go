package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SecurityCenterAutoProvisionResource struct {
}

func TestAccSecurityCenterAutoProvision_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_auto_provisioning", "test")
	r := SecurityCenterAutoProvisionResource{}

	// lintignore:AT001
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.setting("On"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_provision").HasValue("On"),
			),
		},
		data.ImportStep(),
		{
			Config: r.setting("Off"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_provision").HasValue("Off"),
			),
		},
		data.ImportStep(),
	})
}

func (SecurityCenterAutoProvisionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	securityCenterAutoProvisioningName := "default"

	resp, err := clients.SecurityCenter.AutoProvisioningClient.Get(ctx, securityCenterAutoProvisioningName)
	if err != nil {
		return nil, fmt.Errorf("reading Security Center auto provision (%s): %+v", securityCenterAutoProvisioningName, err)
	}

	return utils.Bool(resp.AutoProvisioningSettingProperties != nil), nil
}

func (SecurityCenterAutoProvisionResource) setting(setting string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_auto_provisioning" "test" {
  auto_provision = "%s"
}
`, setting)
}
