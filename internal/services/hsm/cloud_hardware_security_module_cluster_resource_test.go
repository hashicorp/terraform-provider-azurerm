// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hsm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CloudHardwareSecurityModuleClusterResource struct{}

func TestAccCloudHardwareSecurityModuleCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cloud_hardware_security_module_cluster", "test")
	r := CloudHardwareSecurityModuleClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCloudHardwareSecurityModuleCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cloud_hardware_security_module_cluster", "test")
	r := CloudHardwareSecurityModuleClusterResource{}

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

func TestAccCloudHardwareSecurityModuleCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cloud_hardware_security_module_cluster", "test")
	r := CloudHardwareSecurityModuleClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_generated_domain_name_label_scope").HasValue("TenantReuse"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCloudHardwareSecurityModuleCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cloud_hardware_security_module_cluster", "test")
	r := CloudHardwareSecurityModuleClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.environment").HasValue("updated"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCloudHardwareSecurityModuleCluster_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cloud_hardware_security_module_cluster", "test")
	r := CloudHardwareSecurityModuleClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCloudHardwareSecurityModuleCluster_privateEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cloud_hardware_security_module_cluster", "test")
	r := CloudHardwareSecurityModuleClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_endpoint_connections.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			// need another apply to read the private endpoint connection
			Config: r.privateEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_endpoint_connections.#").Exists(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.id").Exists(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.name").Exists(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.type").Exists(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.group_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.private_endpoint.0.id").Exists(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.private_link_service_connection_state.0.status").Exists(),
			),
		},
	})
}

func (CloudHardwareSecurityModuleClusterResource) Exists(ctx context.Context, clientsProvider *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cloudhsmclusters.ParseCloudHsmClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clientsProvider.HSM.CloudHsmClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r CloudHardwareSecurityModuleClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cloud_hardware_security_module_cluster" "import" {
  name                = azurerm_cloud_hardware_security_module_cluster.test.name
  resource_group_name = azurerm_cloud_hardware_security_module_cluster.test.resource_group_name
  location            = azurerm_cloud_hardware_security_module_cluster.test.location
}
`, r.basic(data))
}

func (CloudHardwareSecurityModuleClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cloudhsm-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CloudHardwareSecurityModuleClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cloud_hardware_security_module_cluster" "test" {
  name                = "acctest-hsm-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

}
`, r.template(data), data.RandomString)
}

func (r CloudHardwareSecurityModuleClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_cloud_hardware_security_module_cluster" "test" {
  name                = "acctest-hsm-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  auto_generated_domain_name_label_scope = "TenantReuse"

  tags = {
    environment = "test"
    purpose     = "acceptance-testing"
  }
}
`, r.template(data), data.RandomString, data.RandomString)
}

func (r CloudHardwareSecurityModuleClusterResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_cloud_hardware_security_module_cluster" "test" {
  name                = "acctest-hsm-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    environment = "updated"
    purpose     = "acceptance-testing"
  }
}
`, r.template(data), data.RandomString)
}

func (r CloudHardwareSecurityModuleClusterResource) updateIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctestuai2-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_cloud_hardware_security_module_cluster" "test" {
  name                = "acctest-hsm-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test2.id]
  }

  tags = {
    environment = "updated"
    purpose     = "acceptance-testing"
  }
}
`, r.template(data), data.RandomString)
}

func (r CloudHardwareSecurityModuleClusterResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuai-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_cloud_hardware_security_module_cluster" "test" {
  name                = "acctest-hsm-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

}
`, r.template(data), data.RandomString, data.RandomString)
}

func (r CloudHardwareSecurityModuleClusterResource) privateEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%s"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%s"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctesthsm-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_cloud_hardware_security_module_cluster" "test" {
  name                = "acctest-hsm-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  auto_generated_domain_name_label_scope = "TenantReuse"

  tags = {
    environment = "test"
    purpose     = "private-endpoint-testing"
  }
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-pe-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id

  private_service_connection {
    name                           = "acctest-psc-%s"
    private_connection_resource_id = azurerm_cloud_hardware_security_module_cluster.test.id
    is_manual_connection           = false
    subresource_names              = ["cloudHsm"]
  }
}
`, r.template(data), data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}
