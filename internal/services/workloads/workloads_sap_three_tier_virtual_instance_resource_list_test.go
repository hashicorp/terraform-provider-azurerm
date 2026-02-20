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

func TestAccWorkloadsSAPThreeTierVirtualInstance_list_basic(t *testing.T) {
	r := WorkloadsSapThreeTierVirtualInstanceResource{}
	listResourceAddress := "azurerm_workloads_sap_three_tier_virtual_instance.list"
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_three_tier_virtual_instance", "test1")

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
				Config: r.basicList(data),
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

func (r WorkloadsSapThreeTierVirtualInstanceResource) basicList(data acceptance.TestData) string {
	sapVISNameSuffix := 10 + (data.RandomInteger % 90)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_workloads_sap_three_tier_virtual_instance" "test1" {
  name                        = "X%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "acctestManagedRG1%[3]d"
  app_location                = azurerm_resource_group.app.location
  sap_fqdn                    = "sap1.bpaas.com"

  three_tier_configuration {
    app_resource_group_name = azurerm_resource_group.app.name

    application_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.test.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_D16ds_v4"

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

    central_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.test.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_D16ds_v4"

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

    database_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.test.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_E16ds_v4"

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

resource "azurerm_workloads_sap_three_tier_virtual_instance" "test2" {
  name                        = "Y%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "acctestManagedRG2%[3]d"
  app_location                = azurerm_resource_group.app.location
  sap_fqdn                    = "sap2.bpaas.com"

  three_tier_configuration {
    app_resource_group_name = azurerm_resource_group.app.name

    application_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.test.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_D16ds_v4"

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

    central_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.test.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_D16ds_v4"

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

    database_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.test.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_E16ds_v4"

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

resource "azurerm_workloads_sap_three_tier_virtual_instance" "test3" {
  name                        = "Z%[2]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "acctestManagedRG3%[3]d"
  app_location                = azurerm_resource_group.app.location
  sap_fqdn                    = "sap3.bpaas.com"

  three_tier_configuration {
    app_resource_group_name = azurerm_resource_group.app.name

    application_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.test.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_D16ds_v4"

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

    central_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.test.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_D16ds_v4"

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

    database_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.test.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_E16ds_v4"

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
`, r.template(data), sapVISNameSuffix, data.RandomInteger)
}

func (r WorkloadsSapThreeTierVirtualInstanceResource) basicListQuery() string {
	return `
list "azurerm_workloads_sap_three_tier_virtual_instance" "list" {
  provider = azurerm
  config {}
}
`
}

func (r WorkloadsSapThreeTierVirtualInstanceResource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_workloads_sap_three_tier_virtual_instance" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-sapvis-%[1]d"
  }
}
`, data.RandomInteger)
}
