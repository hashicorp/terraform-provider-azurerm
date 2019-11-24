package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers"
)

func TestAccAzureRMKubernetesClusterNodePool_autoScale(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testAccAzureRMKubernetesClusterNodePool_autoScale(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Disabled
				Config: testAccAzureRMKubernetesClusterNodePool_manualScale(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Enabled
				Config: testAccAzureRMKubernetesClusterNodePool_autoScale(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_autoScaleUpdate(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCount(ri, clientId, clientSecret, location, 1, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCount(ri, clientId, clientSecret, location, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCount(ri, clientId, clientSecret, location, 1, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_availabilityZones(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_availabilityZones(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_errorForAvailabilitySet(t *testing.T) {
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMKubernetesClusterNodePool_availabilitySet(ri, clientId, clientSecret, location),
				ExpectError: regexp.MustCompile("must be a VirtualMachineScaleSet to attach multiple node pools"),
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_multiplePools(t *testing.T) {
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_multiplePools(ri, clientId, clientSecret, location, 3),
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
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScale(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools(t *testing.T) {
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools(ri, clientId, clientSecret, location),
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
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsNodeCount(ri, clientId, clientSecret, location, 1),
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
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsNodeCount(ri, clientId, clientSecret, location, 2),
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
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCount(ri, clientId, clientSecret, location, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// up
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCount(ri, clientId, clientSecret, location, 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			}, {
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// and down
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCount(ri, clientId, clientSecret, location, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(ri, clientId, clientSecret, location, "Standard_F2s_v2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(ri, clientId, clientSecret, location, "Standard_F4s_v2"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_nodePublicIP(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_nodePublicIP(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_nodeTaints(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_nodeTaints(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_manualScale(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMKubernetesClusterNodePool_requiresImport(ri, clientId, clientSecret, location),
				ExpectError: testRequiresImportError("azurerm_kubernetes_cluster_node_pool"),
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_osDiskSizeGB(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_osDiskSizeGB(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_virtualNetworkManual(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_virtualNetworkManual(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_windows(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster_node_pool.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_windows(ri, clientId, clientSecret, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesNodePoolExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKubernetesClusterNodePool_windowsAndLinux(t *testing.T) {
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesClusterNodePool_windowsAndLinux(ri, clientId, clientSecret, location),
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

		client := testAccProvider.Meta().(*ArmClient).Containers.AgentPoolsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

		client := testAccProvider.Meta().(*ArmClient).Containers.AgentPoolsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMKubernetesClusterNodePool_autoScale(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_autoScaleNodeCount(rInt int, clientId, clientSecret, location string, min int, max int) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_availabilitySet(rInt int, clientId, clientSecret, location string) string {
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
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesClusterNodePool_availabilityZones(rInt int, clientId, clientSecret, location string) string {
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
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
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
  vnet_subnet_id = azurerm_subnet.test.id
}
`, rInt, location, rInt, rInt, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesClusterNodePool_manualScale(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsNodeCount(rInt int, clientId, clientSecret, location string, numberOfAgents int) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_manualScaleNodeCount(rInt int, clientId, clientSecret, location string, numberOfAgents int) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_manualScaleVMSku(rInt int, clientId, clientSecret, location, sku string) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_multiplePools(rInt int, clientId, clientSecret, location string, numberOfAgents int) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_nodePublicIP(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_cluster_node_pool" "test" {
  name                   = "internal"
  kubernetes_cluster_id  = azurerm_kubernetes_cluster.test.id
  vm_size                = "Standard_DS2_v2"
  node_count             = 1
  enable_node_public_ip  = true
}
`, template)
}

func testAccAzureRMKubernetesClusterNodePool_nodeTaints(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_requiresImport(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_manualScale(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_osDiskSizeGB(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_template(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateVirtualNetwork(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_virtualNetworkManual(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateVirtualNetwork(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_windows(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateWindows(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_windowsAndLinux(rInt int, clientId, clientSecret, location string) string {
	template := testAccAzureRMKubernetesClusterNodePool_templateWindows(rInt, clientId, clientSecret, location)
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

func testAccAzureRMKubernetesClusterNodePool_template(rInt int, clientId, clientSecret, location string) string {
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
`, rInt, location, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesClusterNodePool_templateVirtualNetwork(rInt int, clientId, clientSecret, location string) string {
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
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesClusterNodePool_templateWindows(rInt int, clientId, clientSecret, location string) string {
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
`, rInt, location, rInt, rInt, clientId, clientSecret)
}
