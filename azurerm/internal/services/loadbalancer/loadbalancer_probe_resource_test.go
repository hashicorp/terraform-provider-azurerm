package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLoadBalancerProbe_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerProbeExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerProbeExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMLoadBalancerProbe_requiresImport),
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerProbeExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerProbe_removal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerProbeIsMissing("azurerm_lb.test", fmt.Sprintf("probe-%d", data.RandomInteger)),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_probe", "test2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_multipleProbes(data, data2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerProbeExists(data.ResourceName),
					testCheckAzureRMLoadBalancerProbeExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data2.ResourceName, "port", "80"),
				),
			},
			data.ImportStep(),
			data2.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerProbe_multipleProbesUpdate(data, data2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerProbeExists(data.ResourceName),
					testCheckAzureRMLoadBalancerProbeExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data2.ResourceName, "port", "8080"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_updateProtocol(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_updateProtocolBefore(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerProbeExists(data.ResourceName),
					resource.TestCheckResourceAttr("azurerm_lb_probe.test", "protocol", "Http"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerProbe_updateProtocolAfter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerProbeExists(data.ResourceName),
					resource.TestCheckResourceAttr("azurerm_lb_probe.test", "protocol", "Tcp"),
				),
			},
		},
	})
}

func testCheckAzureRMLoadBalancerProbeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %q", resourceName)
		}

		id, err := parse.LoadBalancerProbeID(rs.Primary.ID)
		if err != nil {
			return err
		}

		lb, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
		if err != nil {
			if utils.ResponseWasNotFound(lb.Response) {
				return fmt.Errorf("Load Balancer %q (resource group %q) not found for Probe %q", id.LoadBalancerName, id.ResourceGroup, id.ProbeName)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Probe %q", id.LoadBalancerName, id.ResourceGroup, id.ProbeName)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.Probes == nil || len(*props.Probes) == 0 {
			return fmt.Errorf("Probe %q not found in Load Balancer %q (resource group %q)", id.ProbeName, id.LoadBalancerName, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.Probes {
			if v.Name != nil && *v.Name == id.ProbeName {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Probe %q not found in Load Balancer %q (resource group %q)", id.ProbeName, id.LoadBalancerName, id.ResourceGroup)
		}
		return nil
	}
}

func testCheckAzureRMLoadBalancerProbeIsMissing(loadBalancerName string, probeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[loadBalancerName]
		if !ok {
			return fmt.Errorf("not found: %q", loadBalancerName)
		}

		id, err := parse.LoadBalancerID(rs.Primary.ID)
		if err != nil {
			return err
		}

		lb, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if utils.ResponseWasNotFound(lb.Response) {
				return fmt.Errorf("Load Balancer %q (resource group %q) not found while checking for Probe removal", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Probe removal", id.Name, id.ResourceGroup)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.Probes == nil {
			return fmt.Errorf("Probe %q not found in Load Balancer %q (resource group %q)", probeName, id.Name, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.Probes {
			if v.Name != nil && *v.Name == probeName {
				found = true
			}
		}
		if found {
			return fmt.Errorf("Probe %q not removed from Load Balancer %q (resource group %q)", probeName, id.Name, id.ResourceGroup)
		}
		return nil
	}
}

func testAccAzureRMLoadBalancerProbe_basic(data acceptance.TestData) string {
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

func testAccAzureRMLoadBalancerProbe_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLoadBalancerProbe_basic(data)
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

func testAccAzureRMLoadBalancerProbe_removal(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMLoadBalancerProbe_multipleProbes(data, data2 acceptance.TestData) string {
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

func testAccAzureRMLoadBalancerProbe_multipleProbesUpdate(data, data2 acceptance.TestData) string {
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

func testAccAzureRMLoadBalancerProbe_updateProtocolBefore(data acceptance.TestData) string {
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

func testAccAzureRMLoadBalancerProbe_updateProtocolAfter(data acceptance.TestData) string {
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
