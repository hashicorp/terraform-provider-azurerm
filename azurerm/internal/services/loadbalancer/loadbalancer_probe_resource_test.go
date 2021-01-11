package loadbalancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LoadBalancerProbe struct {
}

func TestAccAzureRMLoadBalancerProbe_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	r := LoadBalancerProbe{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerProbe_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	r := LoadBalancerProbe{}

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

func TestAccAzureRMLoadBalancerProbe_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	r := LoadBalancerProbe{}

	data.ResourceTest(t, r, []resource.TestStep{
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleProbes(data, data2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("port").HasValue("80"),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
		{
			Config: r.multipleProbesUpdate(data, data2),
			Check: resource.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.updateProtocolBefore(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protocol").HasValue("Http"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateProtocolAfter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protocol").HasValue("Tcp"),
			),
		},
	})
}

func (r LoadBalancerProbe) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancerProbeID(state.ID)
	if err != nil {
		return nil, err
	}

	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		if utils.ResponseWasNotFound(lb.Response) {
			return nil, fmt.Errorf("Load Balancer %q (resource group %q) not found for Probe %q", id.LoadBalancerName, id.ResourceGroup, id.ProbeName)
		}
		return nil, fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Probe %q", id.LoadBalancerName, id.ResourceGroup, id.ProbeName)
	}
	props := lb.LoadBalancerPropertiesFormat
	if props == nil || props.Probes == nil || len(*props.Probes) == 0 {
		return nil, fmt.Errorf("Probe %q not found in Load Balancer %q (resource group %q)", id.ProbeName, id.LoadBalancerName, id.ResourceGroup)
	}

	found := false
	for _, v := range *props.Probes {
		if v.Name != nil && *v.Name == id.ProbeName {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Probe %q not found in Load Balancer %q (resource group %q)", id.ProbeName, id.LoadBalancerName, id.ResourceGroup)
	}
	return utils.Bool(found), nil
}

func (r LoadBalancerProbe) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancerProbeID(state.ID)
	if err != nil {
		return nil, err
	}

	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Load Balancer %q (Resource Group %q)", id.LoadBalancerName, id.ResourceGroup)
	}
	if lb.LoadBalancerPropertiesFormat == nil {
		return nil, fmt.Errorf("`properties` was nil")
	}
	if lb.LoadBalancerPropertiesFormat.Probes == nil {
		return nil, fmt.Errorf("`properties.Probes` was nil")
	}

	probes := make([]network.Probe, 0)
	for _, probe := range *lb.LoadBalancerPropertiesFormat.Probes {
		if probe.Name == nil || *probe.Name == id.ProbeName {
			continue
		}

		probes = append(probes, probe)
	}
	lb.LoadBalancerPropertiesFormat.Probes = &probes

	future, err := client.LoadBalancers.LoadBalancersClient.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, lb)
	if err != nil {
		return nil, fmt.Errorf("updating Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.LoadBalancers.LoadBalancersClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for update of Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "probe-%d"
  port                = 22
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LoadBalancerProbe) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_probe" "import" {
  name                = azurerm_lb_probe.test.name
  loadbalancer_id     = azurerm_lb_probe.test.loadbalancer_id
  resource_group_name = azurerm_lb_probe.test.resource_group_name
  port                = 22
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "probe-%d"
  port                = 22
}

resource "azurerm_lb_probe" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "probe-%d"
  port                = 80
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "probe-%d"
  port                = 22
}

resource "azurerm_lb_probe" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "probe-%d"
  port                = 8080
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "probe-%d"
  protocol            = "Http"
  request_path        = "/"
  port                = 80
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
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "probe-%d"
  protocol            = "Tcp"
  port                = 80
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
