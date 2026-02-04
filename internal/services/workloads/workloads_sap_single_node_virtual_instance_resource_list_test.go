// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package workloads_test

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

func TestAccWorkloadsSAPSingleNodeVirtualInstance_list_basic(t *testing.T) {
	r := WorkloadsSAPSingleNodeVirtualInstanceResource{}
	listResourceAddress := "azurerm_workloads_sap_single_node_virtual_instance.list"
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_single_node_virtual_instance", "test1")
	sapVISNameSuffix1 := SAPSingleNodeVirtualInstanceNameSuffix()
	sapVISNameSuffix2 := SAPSingleNodeVirtualInstanceNameSuffix()
	sapVISNameSuffix3 := SAPSingleNodeVirtualInstanceNameSuffix()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				VersionConstraint: "=4.1.0",
				Source:            "registry.terraform.io/hashicorp/tls",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data, sapVISNameSuffix1, sapVISNameSuffix2, sapVISNameSuffix3),
			},
			{
				Query:  true,
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) basicList(data acceptance.TestData, sapVISNameSuffix1, sapVISNameSuffix2, sapVISNameSuffix3 int) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_workloads_sap_single_node_virtual_instance" "test1" {
  name                        = "X%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  app_location                = azurerm_resource_group.app.location
  managed_resource_group_name = "acctestManagedRG1%d"
  sap_fqdn                    = "sap1.bpaas.com"

  single_server_configuration {
    app_resource_group_name = azurerm_resource_group.app.name
    subnet_id               = azurerm_subnet.test.id

    virtual_machine_configuration {
      virtual_machine_size = "Standard_E32ds_v4"

      image {
        offer     = "RHEL-SAP-HA"
        publisher = "RedHat"
        sku       = "82sapha-gen2"
        version   = "latest"
      }

      os_profile {
        admin_username  = "testAdmin"
        ssh_private_key = tls_private_key.test.private_key_pem
        ssh_public_key  = data.tls_public_key.test.public_key_openssh
      }
    }
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test,
  ]
}

resource "azurerm_workloads_sap_single_node_virtual_instance" "test2" {
  name                        = "Y%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  app_location                = azurerm_resource_group.app.location
  managed_resource_group_name = "acctestManagedRG2%d"
  sap_fqdn                    = "sap2.bpaas.com"

  single_server_configuration {
    app_resource_group_name = azurerm_resource_group.app.name
    subnet_id               = azurerm_subnet.test.id

    virtual_machine_configuration {
      virtual_machine_size = "Standard_E32ds_v4"

      image {
        offer     = "RHEL-SAP-HA"
        publisher = "RedHat"
        sku       = "82sapha-gen2"
        version   = "latest"
      }

      os_profile {
        admin_username  = "testAdmin"
        ssh_private_key = tls_private_key.test.private_key_pem
        ssh_public_key  = data.tls_public_key.test.public_key_openssh
      }
    }
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test,
  ]
}

resource "azurerm_workloads_sap_single_node_virtual_instance" "test3" {
  name                        = "Z%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  app_location                = azurerm_resource_group.app.location
  managed_resource_group_name = "acctestManagedRG3%d"
  sap_fqdn                    = "sap3.bpaas.com"

  single_server_configuration {
    app_resource_group_name = azurerm_resource_group.app.name
    subnet_id               = azurerm_subnet.test.id

    virtual_machine_configuration {
      virtual_machine_size = "Standard_E32ds_v4"

      image {
        offer     = "RHEL-SAP-HA"
        publisher = "RedHat"
        sku       = "82sapha-gen2"
        version   = "latest"
      }

      os_profile {
        admin_username  = "testAdmin"
        ssh_private_key = tls_private_key.test.private_key_pem
        ssh_public_key  = data.tls_public_key.test.public_key_openssh
      }
    }
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, template, sapVISNameSuffix1, data.RandomInteger, sapVISNameSuffix2, data.RandomInteger, sapVISNameSuffix3, data.RandomInteger)
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) basicListQuery() string {
	return `
list "azurerm_workloads_sap_single_node_virtual_instance" "list" {
  provider = azurerm
  config {}
}
`
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_workloads_sap_single_node_virtual_instance" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-sapvis-%[1]d"
  }
}
`, data.RandomInteger)
}
