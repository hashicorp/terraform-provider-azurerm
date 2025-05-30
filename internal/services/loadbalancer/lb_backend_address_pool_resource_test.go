// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LoadBalancerBackendAddressPool struct{}

// Basic and Standard use different API's for reasons, so we need to test both flows

func TestAccBackendAddressPoolBasicSkuBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolSynchronousModeManual(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.syncModeManual(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolSynchronousModeAutomatic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.syncModeAutomatic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolBasicSkuDisappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicSkuBasic,
			TestResource: r,
		}),
	})
}

func TestAccBackendAddressPoolBasicSkuRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.basicSkuRequiresImport),
	})
}

func TestAccBackendAddressPoolStandardSkuBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardSkuBasicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPoolStandardSkuDisappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.standardSkuBasic,
			TestResource: r,
		}),
	})
}

func TestAccBackendAddressPoolStandardSkuRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.standardSkuRequiresImport),
	})
}

func TestAccBackendAddressPool_GatewaySkuBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gatewaySkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBackendAddressPool_GatewaySkuUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_backend_address_pool", "test")
	r := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gatewaySkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.gatewaySkuComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.gatewaySkuBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LoadBalancerBackendAddressPool) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(lb.HttpResponse) {
			return nil, fmt.Errorf("%s was not found", plbId)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", plbId, err)
	}
	if model := lb.Model; model != nil {
		props := model.Properties
		if props == nil || props.BackendAddressPools == nil || len(*props.BackendAddressPools) == 0 {
			return nil, fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.BackendAddressPoolName, id.LoadBalancerName, id.ResourceGroupName)
		}

		found := false
		for _, v := range *props.BackendAddressPools {
			if v.Name != nil && *v.Name == id.BackendAddressPoolName {
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.BackendAddressPoolName, id.LoadBalancerName, id.ResourceGroupName)
		}
	}

	return pointer.To(true), nil
}

func (r LoadBalancerBackendAddressPool) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("retrieving Load Balancer %q (Resource Group %q)", id.LoadBalancerName, id.ResourceGroupName)
	}
	if lb.Model == nil {
		return nil, fmt.Errorf("`model` was nil")
	}
	if lb.Model.Properties == nil {
		return nil, fmt.Errorf("`properties` was nil")
	}
	if lb.Model.Properties.BackendAddressPools == nil {
		return nil, fmt.Errorf("`properties.BackendAddressPools` was nil")
	}

	backendAddressPools := make([]loadbalancers.BackendAddressPool, 0)
	for _, backendAddressPool := range *lb.Model.Properties.BackendAddressPools {
		if backendAddressPool.Name == nil || *backendAddressPool.Name == id.BackendAddressPoolName {
			continue
		}

		backendAddressPools = append(backendAddressPools, backendAddressPool)
	}
	lb.Model.Properties.BackendAddressPools = &backendAddressPools

	err = client.LoadBalancers.LoadBalancersClient.CreateOrUpdateThenPoll(ctx, plbId, *lb.Model)
	if err != nil {
		return nil, fmt.Errorf("updating Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroupName, err)
	}

	return pointer.To(true), nil
}

func (r LoadBalancerBackendAddressPool) basicSkuBasic(data acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "pool"
  loadbalancer_id = azurerm_lb.test.id
}
`, template)
}

func (r LoadBalancerBackendAddressPool) basicSkuRequiresImport(data acceptance.TestData) string {
	template := r.basicSkuBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "import" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, template)
}

func (r LoadBalancerBackendAddressPool) standardSkuBasic(data acceptance.TestData) string {
	template := r.template(data, "Standard")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "pool"
  loadbalancer_id = azurerm_lb.test.id
}
`, template)
}

func (r LoadBalancerBackendAddressPool) standardSkuBasicUpdate(data acceptance.TestData) string {
	template := r.template(data, "Standard")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_lb_backend_address_pool" "test" {
  name               = "pool"
  loadbalancer_id    = azurerm_lb.test.id
  virtual_network_id = azurerm_virtual_network.test.id
}
`, template)
}

func (r LoadBalancerBackendAddressPool) standardSkuRequiresImport(data acceptance.TestData) string {
	template := r.standardSkuBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "import" {
  name            = azurerm_lb_backend_address_pool.test.name
  loadbalancer_id = azurerm_lb_backend_address_pool.test.loadbalancer_id
}
`, template)
}

func (LoadBalancerBackendAddressPool) template(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
locals {
  number   = %d
  location = %q
  sku      = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${local.number}"
  location = local.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["192.168.0.0/16"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = local.sku
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = local.sku

  frontend_ip_configuration {
    name                 = "feip"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, sku)
}

func (r LoadBalancerBackendAddressPool) gatewaySkuBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "acctest-bap-${local.number}"
  loadbalancer_id = azurerm_lb.test.id
  tunnel_interface {
    identifier = 900
    type       = "Internal"
    protocol   = "VXLAN"
    port       = 15000
  }
}
`, r.templateGateway(data))
}

func (r LoadBalancerBackendAddressPool) gatewaySkuComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "test" {
  name            = "acctest-bap-${local.number}"
  loadbalancer_id = azurerm_lb.test.id
  tunnel_interface {
    identifier = 900
    type       = "Internal"
    protocol   = "VXLAN"
    port       = 15000
  }
  tunnel_interface {
    identifier = 901
    type       = "External"
    protocol   = "VXLAN"
    port       = 15001
  }
  virtual_network_id = azurerm_virtual_network.test.id
}
`, r.templateGateway(data))
}

func (LoadBalancerBackendAddressPool) templateGateway(data acceptance.TestData) string {
	return fmt.Sprintf(`
locals {
  number   = %d
  location = %q
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${local.number}"
  location = local.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-${local.number}"
  resource_group_name  = azurerm_virtual_network.test.resource_group_name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-${local.number}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Gateway"

  frontend_ip_configuration {
    name      = "feip"
    subnet_id = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (LoadBalancerBackendAddressPool) templateSyncMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet-%[1]d"
  resource_group_name  = azurerm_virtual_network.test.resource_group_name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_lb" "test" {
  name                = "acctestlb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name      = "fe-lb"
    subnet_id = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LoadBalancerBackendAddressPool) syncModeManual(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_lb_backend_address_pool" "test" {
  name               = "pool"
  loadbalancer_id    = azurerm_lb.test.id
  synchronous_mode   = "Manual"
  virtual_network_id = azurerm_virtual_network.test.id
}

resource "azurerm_lb_backend_address_pool_address" "test" {
  name                    = "address"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  ip_address              = "10.0.2.6"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"
  priority            = "Regular"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = 1
  }

  boot_diagnostics {
    enabled     = false
    storage_uri = ""
  }

  os_profile {
    computer_name_prefix = "testvm-%[2]d"
    admin_username       = "adminuser"
    admin_password       = "P@ssW0RD7890"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name      = "TestIPConfiguration"
      primary   = true
      subnet_id = azurerm_subnet.test.id
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  depends_on = [azurerm_lb_backend_address_pool_address.test]
}
`, r.templateSyncMode(data), data.RandomInteger)
}

func (r LoadBalancerBackendAddressPool) syncModeAutomatic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%[1]s

resource "azurerm_lb_backend_address_pool" "test" {
  name               = "pool"
  loadbalancer_id    = azurerm_lb.test.id
  synchronous_mode   = "Automatic"
  virtual_network_id = azurerm_virtual_network.test.id
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "acctvmss-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  upgrade_policy_mode = "Manual"
  priority            = "Regular"

  sku {
    name     = "Standard_F2"
    tier     = "Standard"
    capacity = 1
  }

  boot_diagnostics {
    enabled     = false
    storage_uri = ""
  }

  os_profile {
    computer_name_prefix = "testvm-%[2]d"
    admin_username       = "adminuser"
    admin_password       = "P@ssW0RD7890"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  network_profile {
    name    = "TestNetworkProfile"
    primary = true

    ip_configuration {
      name                                   = "TestIPConfiguration"
      primary                                = true
      subnet_id                              = azurerm_subnet.test.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
    }
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}
`, r.templateSyncMode(data), data.RandomInteger)
}
