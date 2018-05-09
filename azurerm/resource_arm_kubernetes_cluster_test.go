package azurerm

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMKubernetesCluster_agentPoolName(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "hi",
			ExpectError: false,
		},
		{
			Input:       "hello",
			ExpectError: false,
		},
		{
			Input:       "hello-world",
			ExpectError: true,
		},
		{
			Input:       "helloworld123",
			ExpectError: true,
		},
		{
			Input:       "hello_world",
			ExpectError: true,
		},
		{
			Input:       "Hello-World",
			ExpectError: true,
		},
		{
			Input:       "20202020",
			ExpectError: true,
		},
		{
			Input:       "h20202020",
			ExpectError: false,
		},
		{
			Input:       "ABC123!@£",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := validateKubernetesClusterAgentPoolName()(tc.Input, "")

		hasError := len(errors) > 0

		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Kubernetes Cluster Agent Pool Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}

func TestAccAzureRMKubernetesCluster_basic(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	config := testAccAzureRMKubernetesCluster_basic(ri, clientId, clientSecret, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "kube_config.0.password"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_addAgent(t *testing.T) {
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	initConfig := testAccAzureRMKubernetesCluster_basic(ri, clientId, clientSecret, testLocation())
	addAgentConfig := testAccAzureRMKubernetesCluster_addAgent(ri, clientId, clientSecret, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: initConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists("azurerm_kubernetes_cluster.test"),
				),
			},
			{
				Config: addAgentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("azurerm_kubernetes_cluster.test", "agent_pool_profile.0.count", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_upgradeConfig(t *testing.T) {
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	initConfig := testAccAzureRMKubernetesCluster_basic(ri, clientId, clientSecret, testLocation())
	upgradeConfig := testAccAzureRMKubernetesCluster_upgrade(ri, clientId, clientSecret, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: initConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists("azurerm_kubernetes_cluster.test"),
				),
			},
			{
				Config: upgradeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("azurerm_kubernetes_cluster.test", "kubernetes_version", "1.8.1"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	config := testAccAzureRMKubernetesCluster_internalNetwork(ri, clientId, clientSecret, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists("azurerm_kubernetes_cluster.test"),
				),
			},
		},
	})
}

func testAccAzureRMKubernetesCluster_basic(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "1.7.7"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name    = "default"
    count   = "1"
    vm_size = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_addAgent(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "1.7.7"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name    = "default"
    count   = "2"
    vm_size = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_internalNetwork(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags {
    environment = "Testing"
  }
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                   = "acctestaks%d"
  location               = "${azurerm_resource_group.test.location}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  dns_prefix             = "acctestaks%d"
  kubernetes_version     = "1.7.7"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name           = "default"
    count          = "2"
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = "${azurerm_subnet.test.id}"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_upgrade(rInt int, clientId string, clientSecret string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                   = "acctestaks%d"
  location               = "${azurerm_resource_group.test.location}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  dns_prefix             = "acctestaks%d"
  kubernetes_version     = "1.8.1"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  agent_pool_profile {
    name       = "default"
    count      = "1"
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}
`, rInt, location, rInt, rInt, rInt, clientId, clientSecret)
}

func testCheckAzureRMKubernetesClusterExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for AKS managed cluster instance: %s", resourceName)
		}

		client := testAccProvider.Meta().(*ArmClient).kubernetesClustersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		aks, err := client.Get(ctx, resourceGroup, resourceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on kubernetesClustersClient: %+v", err)
		}

		if aks.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: AKS managed cluster instance %q (resource group: %q) does not exist", resourceName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMKubernetesClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).kubernetesClustersClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kubernetes_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("AKS managed cluster instance still exists:\n%#v", resp)
		}
	}

	return nil
}
