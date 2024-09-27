// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedDiskSASTokenResource struct{}

func TestAccManagedDiskSASToken_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_disk_sas_token", "test")
	r := ManagedDiskSASTokenResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t ManagedDiskSASTokenResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseManagedDiskID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.DisksClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Disk Export status %q", id.String())
	}

	if string(*resp.Model.Properties.DiskState) != "ActiveSAS" {
		return nil, fmt.Errorf("Disk SAS token %s (resource group %s): %s", id.DiskName, id.ResourceGroupName, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ManagedDiskSASTokenResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-revokedisk-%d"
  location = "%s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestsads%s"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
}

resource "azurerm_managed_disk_sas_token" "test" {
  managed_disk_id     = azurerm_managed_disk.test.id
  duration_in_seconds = 300
  access_level        = "Read"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
