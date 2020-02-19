package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMHDInsightRServerCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_rserver_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy(data.ResourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightRServerCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
			data.ImportStep("roles.0.head_node.0.password",
				"roles.0.head_node.0.vm_size",
				"roles.0.worker_node.0.password",
				"roles.0.worker_node.0.vm_size",
				"roles.0.zookeeper_node.0.password",
				"roles.0.zookeeper_node.0.vm_size",
				"roles.0.edge_node.0.password",
				"roles.0.edge_node.0.vm_size",
				"storage_account"),
		},
	})
}

func TestAccAzureRMHDInsightRServerCluster_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_hdinsight_rserver_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy(data.ResourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightRServerCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMHDInsightRServerCluster_requiresImport),
		},
	})
}

func TestAccAzureRMHDInsightRServerCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_rserver_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy(data.ResourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightRServerCluster_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
			data.ImportStep("roles.0.head_node.0.password",
				"roles.0.head_node.0.vm_size",
				"roles.0.worker_node.0.password",
				"roles.0.worker_node.0.vm_size",
				"roles.0.zookeeper_node.0.password",
				"roles.0.zookeeper_node.0.vm_size",
				"roles.0.edge_node.0.password",
				"roles.0.edge_node.0.vm_size",
				"storage_account"),
			{
				Config: testAccAzureRMHDInsightRServerCluster_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
			data.ImportStep("roles.0.head_node.0.password",
				"roles.0.head_node.0.vm_size",
				"roles.0.worker_node.0.password",
				"roles.0.worker_node.0.vm_size",
				"roles.0.zookeeper_node.0.password",
				"roles.0.zookeeper_node.0.vm_size",
				"roles.0.edge_node.0.password",
				"roles.0.edge_node.0.vm_size",
				"storage_account"),
		},
	})
}

func TestAccAzureRMHDInsightRServerCluster_sshKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_rserver_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy(data.ResourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightRServerCluster_sshKeys(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
			data.ImportStep("storage_account",
				"roles.0.head_node.0.ssh_keys",
				"roles.0.head_node.0.vm_size",
				"roles.0.worker_node.0.ssh_keys",
				"roles.0.worker_node.0.vm_size",
				"roles.0.zookeeper_node.0.ssh_keys",
				"roles.0.zookeeper_node.0.vm_size",
				"roles.0.edge_node.0.ssh_keys",
				"roles.0.edge_node.0.vm_size"),
		},
	})
}

func TestAccAzureRMHDInsightRServerCluster_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_rserver_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy(data.ResourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightRServerCluster_virtualNetwork(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
			data.ImportStep("roles.0.head_node.0.password",
				"roles.0.head_node.0.vm_size",
				"roles.0.worker_node.0.password",
				"roles.0.worker_node.0.vm_size",
				"roles.0.zookeeper_node.0.password",
				"roles.0.zookeeper_node.0.vm_size",
				"roles.0.edge_node.0.password",
				"roles.0.edge_node.0.vm_size",
				"storage_account"),
		},
	})
}

func TestAccAzureRMHDInsightRServerCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_rserver_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy(data.ResourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightRServerCluster_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
			data.ImportStep("roles.0.head_node.0.password",
				"roles.0.head_node.0.vm_size",
				"roles.0.worker_node.0.password",
				"roles.0.worker_node.0.vm_size",
				"roles.0.zookeeper_node.0.password",
				"roles.0.zookeeper_node.0.vm_size",
				"roles.0.edge_node.0.password",
				"roles.0.edge_node.0.vm_size",
				"storage_account"),
		},
	})
}

func testAccAzureRMHDInsightRServerCluster_basic(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightRServerCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_rserver_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"
  rstudio             = true

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = "${azurerm_storage_container.test.id}"
    storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
    is_default           = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 2
    }

    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    edge_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMHDInsightRServerCluster_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightRServerCluster_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_rserver_cluster" "import" {
  name                = "${azurerm_hdinsight_rserver_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_rserver_cluster.test.resource_group_name}"
  location            = "${azurerm_hdinsight_rserver_cluster.test.location}"
  cluster_version     = "${azurerm_hdinsight_rserver_cluster.test.cluster_version}"
  tier                = "${azurerm_hdinsight_rserver_cluster.test.tier}"
  gateway             = "${azurerm_hdinsight_rserver_cluster.test.gateway}"
  storage_account     = "${azurerm_hdinsight_rserver_cluster.test.storage_account}"
  roles               = "${azurerm_hdinsight_rserver_cluster.test.roles}"
}
`, template)
}

func testAccAzureRMHDInsightRServerCluster_sshKeys(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightRServerCluster_template(data)
	return fmt.Sprintf(`
%s

variable "ssh_key" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
}

resource "azurerm_hdinsight_rserver_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"
  rstudio             = true

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = "${azurerm_storage_container.test.id}"
    storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
    is_default           = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      ssh_keys = ["${var.ssh_key}"]
    }

    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      ssh_keys              = ["${var.ssh_key}"]
      target_instance_count = 3
    }

    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      ssh_keys = ["${var.ssh_key}"]
    }

    edge_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      ssh_keys = ["${var.ssh_key}"]
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMHDInsightRServerCluster_updated(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightRServerCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_rserver_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"
  rstudio             = true

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = "${azurerm_storage_container.test.id}"
    storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
    is_default           = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 5
    }

    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    edge_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }

  tags = {
    Hello = "World"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMHDInsightRServerCluster_virtualNetwork(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightRServerCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_hdinsight_rserver_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"
  rstudio             = true

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = "${azurerm_storage_container.test.id}"
    storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
    is_default           = true
  }

  roles {
    head_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = "${azurerm_subnet.test.id}"
      virtual_network_id = "${azurerm_virtual_network.test.id}"
    }

    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
      subnet_id             = "${azurerm_subnet.test.id}"
      virtual_network_id    = "${azurerm_virtual_network.test.id}"
    }

    zookeeper_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = "${azurerm_subnet.test.id}"
      virtual_network_id = "${azurerm_virtual_network.test.id}"
    }

    edge_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = "${azurerm_subnet.test.id}"
      virtual_network_id = "${azurerm_virtual_network.test.id}"
    }
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMHDInsightRServerCluster_complete(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightRServerCluster_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_hdinsight_rserver_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"
  rstudio             = true

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = "${azurerm_storage_container.test.id}"
    storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
    is_default           = true
  }

  roles {
    head_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = "${azurerm_subnet.test.id}"
      virtual_network_id = "${azurerm_virtual_network.test.id}"
    }

    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
      subnet_id             = "${azurerm_subnet.test.id}"
      virtual_network_id    = "${azurerm_virtual_network.test.id}"
    }

    zookeeper_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = "${azurerm_subnet.test.id}"
      virtual_network_id = "${azurerm_virtual_network.test.id}"
    }

    edge_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = "${azurerm_subnet.test.id}"
      virtual_network_id = "${azurerm_virtual_network.test.id}"
    }
  }

  tags = {
    Hello = "World"
  }
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMHDInsightRServerCluster_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
