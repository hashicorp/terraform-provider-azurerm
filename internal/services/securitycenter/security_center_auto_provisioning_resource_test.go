// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityCenterAutoProvisionResource struct{}

func TestAccSecurityCenterAutoProvision_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_auto_provisioning", "test")
	r := SecurityCenterAutoProvisionResource{}

	//lintignore:AT001
	data.ResourceTestSkipCheckDestroyed(t, []acceptance.TestStep{
		{
			Config: r.setting("On"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_provision").HasValue("On"),
			),
		},
		data.ImportStep(),
		{
			Config: r.setting("Off"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_provision").HasValue("Off"),
			),
		},
		data.ImportStep(),
	})
}

func (SecurityCenterAutoProvisionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AutoProvisioningSettingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.AutoProvisioningClient.Get(ctx, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving auto-provisioning setting for %s: %+v", *id, err)
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
