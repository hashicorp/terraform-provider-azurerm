package tests

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers"
)

func TestAccAzureRMKubernetesClusterNodePool_autoScale(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_autoScale(t)
}

func testAccAzureRMKubernetesClusterNodePool_autoScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Disabled
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// Enabled
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_autoScaleUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_autoScaleUpdate(t)
}

func testAccAzureRMKubernetesClusterNodePool_autoScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCountConfig(data, clientId, clientSecret, data.Locations.Primary, 1, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCountConfig(data, clientId, clientSecret, data.Locations.Primary, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCountConfig(data, clientId, clientSecret, data.Locations.Primary, 1, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_availabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_availabilityZones(t)
}

func testAccAzureRMKubernetesClusterNodePool_availabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_availabilityZonesConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_errorForAvailabilitySet(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_errorForAvailabilitySet(t)
}

func testAccAzureRMKubernetesClusterNodePool_errorForAvailabilitySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMKubernetesClusterNodePool_availabilitySetConfig(data, clientId, clientSecret),
				ExpectError: regexp.MustCompile("must be a VirtualMachineScaleSet to attach multiple node pools"),
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_multiplePools(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_multiplePools(t)
}

func testAccAzureRMKubernetesClusterNodePool_multiplePools(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_multiplePoolsConfig(data, clientId, clientSecret, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.autoscale"),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.manual"),
				),
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.autoscale",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.manual",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScale(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScale(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.first"),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.second"),
				),
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.first",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.second",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsNodeCountConfig(data, clientId, clientSecret, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.first"),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.second"),
				),
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.first",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.second",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsNodeCountConfig(data, clientId, clientSecret, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.first"),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.second"),
				),
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.first",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.second",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleUpdate(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScaleUpdate(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCountConfig(data, clientId, clientSecret, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				// up
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCountConfig(data, clientId, clientSecret, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			}, data.ImportStep(),
			{
				// and down
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCountConfig(data, clientId, clientSecret, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(t)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleVMSkuConfig(data, clientId, clientSecret, "Standard_F2s_v2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleVMSkuConfig(data, clientId, clientSecret, "Standard_F4s_v2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_nodePublicIP(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_nodePublicIP(t)
}

func testAccAzureRMKubernetesClusterNodePool_nodePublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_nodePublicIPConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_nodeTaints(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_nodeTaints(t)
}

func testAccAzureRMKubernetesClusterNodePool_nodeTaints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_nodeTaintsConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_requiresImport(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_requiresImport(t)
}

func testAccAzureRMKubernetesClusterNodePool_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMKubernetesClusterNodePool_requiresImportConfig(data, clientId, clientSecret),
				ExpectError: acceptance.RequiresImportError("azurerm_kubernetes_cluster_node_pool"),
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_osDiskSizeGB(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_osDiskSizeGB(t)
}

func testAccAzureRMKubernetesClusterNodePool_osDiskSizeGB(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_osDiskSizeGBConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic(t)
}

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomaticConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_virtualNetworkManual(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_virtualNetworkManual(t)
}

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkManual(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_virtualNetworkManualConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_windows(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_windows(t)
}

func testAccAzureRMKubernetesClusterNodePool_windows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_windowsConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_windowsAndLinux(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesClusterNodePool_windowsAndLinux(t)
}

func testAccAzureRMKubernetesClusterNodePool_windowsAndLinux(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_windowsAndLinuxConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.linux"),
					testCheckAzureRMKubernetesNodePoolExists("azurerm_kubernetes_cluster_node_pool.windows"),
				),
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.linux",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "azurerm_kubernetes_cluster_node_pool.windows",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMKubernetesClusterNodePoolDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.AgentPoolsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kubernetes_cluster_node_pool" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		kubernetesClusterId := rs.Primary.Attributes["kubernetes_cluster_id"]
		parsedK8sId, err := containers.ParseKubernetesClusterID(kubernetesClusterId)
		if err != nil {
			return fmt.Errorf("Error parsing kubernetes cluster id: %+v", err)
		}

		resp, err := client.Get(ctx, parsedK8sId.ResourceGroup, parsedK8sId.Name, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Managed Kubernetes Cluster still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMKubernetesNodePoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.AgentPoolsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		kubernetesClusterId := rs.Primary.Attributes["kubernetes_cluster_id"]
		parsedK8sId, err := containers.ParseKubernetesClusterID(kubernetesClusterId)
		if err != nil {
			return fmt.Errorf("Error parsing kubernetes cluster id: %+v", err)
		}

		agent_pool, err := client.Get(ctx, parsedK8sId.ResourceGroup, parsedK8sId.Name, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on kubernetesClustersClient: %+v", err)
		}

		if agent_pool.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Node Pool %q (Kubernetes Cluster %q / Resource Group: %q) does not exist", name, parsedK8sId.Name, parsedK8sId.ResourceGroup)
		}

		return nil
	}
}

func testAccAzureRMKubernetesClusterNodePool_autoScaleConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = 1
  max_count             = 3
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCountConfig(data acceptance.TestData, clientId, clientSecret, location string, min int, max int) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = %d
  max_count             = %d
}
`, template, min, max)
}

func testAccAzureRMKubernetesClusterNodePool_availabilitySetConfig(data acceptance.TestData, clientId, clientSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    type       = "AvailabilitySet"
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesClusterNodePool_availabilityZonesConfig(data acceptance.TestData, clientId, clientSecret string) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name           = "default"
    node_count     = 1
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "Standard"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  availability_zones    = ["1"]
  vnet_subnet_id        = azurerm_subnet.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "first" {
  name                  = "first"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}

resource "azurerm_kubernetes_cluster_node_pool" "second" {
  name                  = "second"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_F2s_v2"
  node_count            = 1
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsNodeCountConfig(data acceptance.TestData, clientId, clientSecret string, numberOfAgents int) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "first" {
  name                  = "first"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = %d
}

resource "azurerm_kubernetes_cluster_node_pool" "second" {
  name                  = "second"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_F2s_v2"
  node_count            = %d
}
`, template, numberOfAgents, numberOfAgents)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCountConfig(data acceptance.TestData, clientId, clientSecret string, numberOfAgents int) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = %d
}
`, template, numberOfAgents)
}

func testAccAzureRMKubernetesClusterNodePool_manualScaleVMSkuConfig(data acceptance.TestData, clientId, clientSecret, sku string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "%s"
  node_count            = 1
}
`, template, sku)
}

func testAccAzureRMKubernetesClusterNodePool_multiplePoolsConfig(data acceptance.TestData, clientId, clientSecret string, numberOfAgents int) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "autoscale" {
  name                  = "autoscale"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = 1
  max_count             = 3
}

resource "azurerm_kubernetes_cluster_node_pool" "manual" {
  name                  = "manual"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_F2s_v2"
  node_count            = %d
}
`, template, numberOfAgents)
}

func testAccAzureRMKubernetesClusterNodePool_nodePublicIPConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  enable_node_public_ip = true
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_nodeTaintsConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  node_taints = [
    "key=value:NoSchedule"
  ]
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_requiresImportConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_manualScaleConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "import" {
  name                  = azurerm_kubernetes_cluster_node_pool.test.name
  kubernetes_cluster_id = azurerm_kubernetes_cluster_node_pool.test.kubernetes_cluster_id
  vm_size               = azurerm_kubernetes_cluster_node_pool.test.vm_size
  node_count            = azurerm_kubernetes_cluster_node_pool.test.node_count
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_osDiskSizeGBConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  os_disk_size_gb       = 100
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomaticConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateVirtualNetworkConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  enable_auto_scaling   = true
  min_count             = 1
  max_count             = 3
  vnet_subnet_id        = azurerm_subnet.test.id
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkManualConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateVirtualNetworkConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  vnet_subnet_id        = azurerm_subnet.test.id
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_windowsConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateWindowsConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                  = "windoz"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  os_type               = "Windows"
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_windowsAndLinuxConfig(data acceptance.TestData, clientId, clientSecret string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateWindowsConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "linux" {
  name                  = "linux"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}

resource "azurerm_kubernetes_cluster_node_pool" "windows" {
  name                  = "windoz"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.test.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
  os_type               = "Windows"
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_templateConfig(data acceptance.TestData, clientId, clientSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesClusterNodePool_templateVirtualNetworkConfig(data acceptance.TestData, clientId, clientSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name                   = "akc-route-%d"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"

  # TODO: remove in 2.0
  lifecycle {
    ignore_changes = ["route_table_id"]
  }
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name           = "default"
    node_count     = 1
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesClusterNodePool_templateWindowsConfig(data acceptance.TestData, clientId, clientSecret string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  windows_profile {
    admin_username = "azureuser"
    admin_password = "P@55W0rd1234!"
  }

  network_profile {
    network_plugin     = "azure"
    network_policy     = "azure"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}
