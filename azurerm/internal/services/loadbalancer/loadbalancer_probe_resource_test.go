package loadbalancer_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer"
)

func TestAccAzureRMLoadBalancerProbe_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")

	var lb network.LoadBalancer
	probeName := fmt.Sprintf("probe-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	probeId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/probes/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, probeName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(data, probeName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_probe.test", "id", probeId),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")

	var lb network.LoadBalancer
	probeName := fmt.Sprintf("probe-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	probeId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/probes/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, probeName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(data, probeName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_probe.test", "id", probeId),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerProbe_requiresImport(data, probeName),
				ExpectError: acceptance.RequiresImportError("azurerm_lb_probe"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")

	var lb network.LoadBalancer
	probeName := fmt.Sprintf("probe-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(data, probeName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerProbe_removal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeNotExists(probeName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_update(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_probe", "test2")

	var lb network.LoadBalancer
	probeName := fmt.Sprintf("probe-%d", data1.RandomInteger)
	probe2Name := fmt.Sprintf("probe-%d", data2.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_multipleProbes(data1, probeName, probe2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					testCheckAzureRMLoadBalancerProbeExists(probe2Name, &lb),
					resource.TestCheckResourceAttr(data2.ResourceName, "port", "80"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerProbe_multipleProbesUpdate(data1, probeName, probe2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					testCheckAzureRMLoadBalancerProbeExists(probe2Name, &lb),
					resource.TestCheckResourceAttr(data2.ResourceName, "port", "8080"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_updateProtocol(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")

	var lb network.LoadBalancer
	probeName := fmt.Sprintf("probe-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_updateProtocolBefore(data, probeName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_probe.test", "protocol", "Http"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerProbe_updateProtocolAfter(data, probeName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_probe.test", "protocol", "Tcp"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerProbe_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_probe", "test")

	var lb network.LoadBalancer
	probeName := fmt.Sprintf("probe-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerProbe_basic(data, probeName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
					testCheckAzureRMLoadBalancerProbeDisappears(probeName, &lb),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLoadBalancerProbeExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := loadbalancer.FindLoadBalancerProbeByName(lb, natRuleName)
		if !exists {
			return fmt.Errorf("A Probe with name %q cannot be found.", natRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerProbeNotExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := loadbalancer.FindLoadBalancerProbeByName(lb, natRuleName)
		if exists {
			return fmt.Errorf("A Probe with name %q has been found.", natRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerProbeDisappears(addressPoolName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		_, i, exists := loadbalancer.FindLoadBalancerProbeByName(lb, addressPoolName)
		if !exists {
			return fmt.Errorf("A Probe with name %q cannot be found.", addressPoolName)
		}

		currentProbes := *lb.LoadBalancerPropertiesFormat.Probes
		probes := append(currentProbes[:i], currentProbes[i+1:]...)
		lb.LoadBalancerPropertiesFormat.Probes = &probes

		id, err := azure.ParseAzureResourceID(*lb.ID)
		if err != nil {
			return err
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, *lb.Name, *lb)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating LoadBalancer: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for completion for LoadBalancer: %+v", err)
		}

		_, err = client.Get(ctx, id.ResourceGroup, *lb.Name, "")
		return err
	}
}

func testAccAzureRMLoadBalancerProbe_basic(data acceptance.TestData, probeName string) string {
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
  name                = "%s"
  port                = 22
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, probeName)
}

func testAccAzureRMLoadBalancerProbe_requiresImport(data acceptance.TestData, name string) string {
	template := testAccAzureRMLoadBalancerProbe_basic(data, name)
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

func testAccAzureRMLoadBalancerProbe_multipleProbes(data acceptance.TestData, probeName, probe2Name string) string {
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
  name                = "%s"
  port                = 22
}

resource "azurerm_lb_probe" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "%s"
  port                = 80
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, probeName, probe2Name)
}

func testAccAzureRMLoadBalancerProbe_multipleProbesUpdate(data acceptance.TestData, probeName, probe2Name string) string {
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
  name                = "%s"
  port                = 22
}

resource "azurerm_lb_probe" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "%s"
  port                = 8080
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, probeName, probe2Name)
}

func testAccAzureRMLoadBalancerProbe_updateProtocolBefore(data acceptance.TestData, probeName string) string {
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
  name                = "%s"
  protocol            = "Http"
  request_path        = "/"
  port                = 80
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, probeName)
}

func testAccAzureRMLoadBalancerProbe_updateProtocolAfter(data acceptance.TestData, probeName string) string {
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
  name                = "%s"
  protocol            = "Tcp"
  port                = 80
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, probeName)
}
