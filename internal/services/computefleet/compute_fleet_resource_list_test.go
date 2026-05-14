// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package computefleet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccComputeFleet_list_basic(t *testing.T) {
	r := ComputeFleetResource{}
	listResourceAddress := "azurerm_compute_fleet.list"

	data := acceptance.BuildTestData(t, "azurerm_compute_fleet", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r ComputeFleetResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_compute_fleet" "test" {
  count = 3

  name                = "acctest-fleet-${count.index}-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  spot_capacity {
    maintain_capacity_enabled = false
    target_capacity           = 1
  }

  virtual_machine_sizes_profile {
    name = "Standard_F1alds_v7"
  }

  virtual_machine_profile {
    source_image_reference {
      offer     = "0001-com-ubuntu-server-focal"
      publisher = "canonical"
      sku       = "20_04-lts-gen2"
      version   = "latest"
    }

    os_profile {
      linux_configuration {
        admin_username                  = local.admin_username
        admin_password                  = local.admin_password
        password_authentication_enabled = true
      }
    }

    network_interface {
      name = "acctest-nic-${count.index}-%[2]d"

      ip_configuration {
        name      = "acctest-ipconfig-${count.index}-%[2]d"
        subnet_id = azurerm_subnet.test.id
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (ComputeFleetResource) basicQuery() string {
	return `
list "azurerm_compute_fleet" "list" {
  provider = azurerm
  config {}
}
`
}

func (ComputeFleetResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_compute_fleet" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
