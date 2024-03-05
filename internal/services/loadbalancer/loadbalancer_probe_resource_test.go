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

type LoadBalancerProbe struct{}

func TestAccAzureRMLoadBalancerProbe_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	r := LoadBalancerProbe{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerProbe_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	r := LoadBalancerProbe{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerProbe_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	r := LoadBalancerProbe{}

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

func TestAccAzureRMLoadBalancerProbe_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	r := LoadBalancerProbe{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccAzureRMLoadBalancerProbe_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_probe", "test2")
	r := LoadBalancerProbe{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleProbes(data, data2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("port").HasValue("80"),
				check.That(data2.ResourceName).Key("probe_threshold").HasValue("3"),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
		{
			Config: r.multipleProbesUpdate(data, data2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("port").HasValue("8080"),
			),
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_updateProtocol(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	r := LoadBalancerProbe{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateProtocolBefore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protocol").HasValue("Http"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateProtocolAfter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protocol").HasValue("Tcp"),
			),
		},
	})
}

func (r LoadBalancerProbe) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseProbeID(state.ID)
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
	props := lb.Model
	if props == nil || props.Properties.Probes == nil || len(*props.Properties.Probes) == 0 {
		return nil, fmt.Errorf("Probe %q not found in Load Balancer %q (resource group %q)", id.ProbeName, id.LoadBalancerName, id.ResourceGroupName)
	}

	found := false
	for _, v := range *props.Properties.Probes {
		if v.Name != nil && *v.Name == id.ProbeName {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Probe %q not found in Load Balancer %q (resource group %q)", id.ProbeName, id.LoadBalancerName, id.ResourceGroupName)
	}
	return pointer.To(found), nil
}

func (r LoadBalancerProbe) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := loadbalancers.ParseProbeID(state.ID)
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
	if lb.Model.Properties.Probes == nil {
		return nil, fmt.Errorf("`properties.Probes` was nil")
	}

	probes := make([]loadbalancers.Probe, 0)
	for _, probe := range *lb.Model.Properties.Probes {
		if probe.Name == nil || *probe.Name == id.ProbeName {
			continue
		}

		probes = append(probes, probe)
	}
	lb.Model.Properties.Probes = &probes

	err = client.LoadBalancers.LoadBalancersClient.CreateOrUpdateThenPoll(ctx, plbId, *lb.Model)
	if err != nil {
		return nil, fmt.Errorf("updating %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (r LoadBalancerProbe) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_probe" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "probe-%d"
  port            = 22
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerProbe) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "one-%[1]d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_probe" "test" {
  loadbalancer_id     = azurerm_lb.test.id
  name                = "probe-%[1]d"
  port                = 22
  interval_in_seconds = 5
  number_of_probes    = 2
  probe_threshold     = 2
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LoadBalancerProbe) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_probe" "import" {
  name            = azurerm_lb_probe.test.name
  loadbalancer_id = azurerm_lb_probe.test.loadbalancer_id
  port            = 22
}
`, template)
}

func (r LoadBalancerProbe) multipleProbes(data, data2 acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_probe" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "probe-%d"
  port            = 22
  probe_threshold = 2
}

resource "azurerm_lb_probe" "test2" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "probe-%d"
  port            = 80
  probe_threshold = 3
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data2.RandomInteger)
}

func (r LoadBalancerProbe) multipleProbesUpdate(data, data2 acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_probe" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "probe-%d"
  port            = 22
}

resource "azurerm_lb_probe" "test2" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "probe-%d"
  port            = 8080
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data2.RandomInteger)
}

func (r LoadBalancerProbe) updateProtocolBefore(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_probe" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "probe-%d"
  protocol        = "Http"
  request_path    = "/"
  port            = 80
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerProbe) updateProtocolAfter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_probe" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "probe-%d"
  protocol        = "Tcp"
  port            = 80
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
