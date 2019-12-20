package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMHDInsightInteractiveQueryCluster_basic(t *testing.T) {
	resourceName := "azurerm_hdinsight_interactive_query_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy("azurerm_hdinsight_interactive_query_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightInteractiveQueryCluster_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"roles.0.head_node.0.password",
					"roles.0.head_node.0.vm_size",
					"roles.0.worker_node.0.password",
					"roles.0.worker_node.0.vm_size",
					"roles.0.zookeeper_node.0.password",
					"roles.0.zookeeper_node.0.vm_size",
					"storage_account",
				},
			},
		},
	})
}

func TestAccAzureRMHDInsightInteractiveQueryCluster_gen2basic(t *testing.T) {
	resourceName := "azurerm_hdinsight_interactive_query_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy("azurerm_hdinsight_interactive_query_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightInteractiveQueryCluster_gen2basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"roles.0.head_node.0.password",
					"roles.0.head_node.0.vm_size",
					"roles.0.worker_node.0.password",
					"roles.0.worker_node.0.vm_size",
					"roles.0.zookeeper_node.0.password",
					"roles.0.zookeeper_node.0.vm_size",
					"storage_account",
				},
			},
		},
	})
}

func TestAccAzureRMHDInsightInteractiveQueryCluster_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_hdinsight_interactive_query_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy("azurerm_hdinsight_interactive_query_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightInteractiveQueryCluster_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_endpoint"),
				),
			},
			{
				Config:      testAccAzureRMHDInsightInteractiveQueryCluster_requiresImport(ri, rs, location),
				ExpectError: acceptance.RequiresImportError("azurerm_hdinsight_interactive_query_cluster"),
			},
		},
	})
}

func TestAccAzureRMHDInsightInteractiveQueryCluster_update(t *testing.T) {
	resourceName := "azurerm_hdinsight_interactive_query_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy("azurerm_hdinsight_interactive_query_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightInteractiveQueryCluster_basic(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"roles.0.head_node.0.password",
					"roles.0.head_node.0.vm_size",
					"roles.0.worker_node.0.password",
					"roles.0.worker_node.0.vm_size",
					"roles.0.zookeeper_node.0.password",
					"roles.0.zookeeper_node.0.vm_size",
					"storage_account",
				},
			},
			{
				Config: testAccAzureRMHDInsightInteractiveQueryCluster_updated(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"roles.0.head_node.0.password",
					"roles.0.head_node.0.vm_size",
					"roles.0.worker_node.0.password",
					"roles.0.worker_node.0.vm_size",
					"roles.0.zookeeper_node.0.password",
					"roles.0.zookeeper_node.0.vm_size",
					"storage_account",
				},
			},
		},
	})
}

func TestAccAzureRMHDInsightInteractiveQueryCluster_sshKeys(t *testing.T) {
	resourceName := "azurerm_hdinsight_interactive_query_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy("azurerm_hdinsight_interactive_query_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightInteractiveQueryCluster_sshKeys(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"storage_account",
					"roles.0.head_node.0.ssh_keys",
					"roles.0.head_node.0.vm_size",
					"roles.0.worker_node.0.ssh_keys",
					"roles.0.worker_node.0.vm_size",
					"roles.0.zookeeper_node.0.ssh_keys",
					"roles.0.zookeeper_node.0.vm_size",
				},
			},
		},
	})
}

func TestAccAzureRMHDInsightInteractiveQueryCluster_virtualNetwork(t *testing.T) {
	resourceName := "azurerm_hdinsight_interactive_query_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy("azurerm_hdinsight_interactive_query_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightInteractiveQueryCluster_virtualNetwork(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"roles.0.head_node.0.password",
					"roles.0.head_node.0.vm_size",
					"roles.0.worker_node.0.password",
					"roles.0.worker_node.0.vm_size",
					"roles.0.zookeeper_node.0.password",
					"roles.0.zookeeper_node.0.vm_size",
					"storage_account",
				},
			},
		},
	})
}

func TestAccAzureRMHDInsightInteractiveQueryCluster_complete(t *testing.T) {
	resourceName := "azurerm_hdinsight_interactive_query_cluster.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy("azurerm_hdinsight_interactive_query_cluster"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightInteractiveQueryCluster_complete(ri, rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "ssh_endpoint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"roles.0.head_node.0.password",
					"roles.0.head_node.0.vm_size",
					"roles.0.worker_node.0.password",
					"roles.0.worker_node.0.vm_size",
					"roles.0.zookeeper_node.0.password",
					"roles.0.zookeeper_node.0.vm_size",
					"storage_account",
				},
			},
		},
	})
}

func testAccAzureRMHDInsightInteractiveQueryCluster_basic(rInt int, rString string, location string) string {
	template := testAccAzureRMHDInsightInteractiveQueryCluster_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_interactive_query_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"

  component_version {
    interactive_hive = "2.1"
  }

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
      vm_size  = "Standard_D13_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_D14_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 2
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, template, rInt)
}

func testAccAzureRMHDInsightInteractiveQueryCluster_gen2basic(rInt int, rString string, location string) string {
	template := testAccAzureRMHDInsightInteractiveQueryCluster_gen2template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_interactive_query_cluster" "test" {
  depends_on = [azurerm_role_assignment.test]

  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"

  component_version {
    interactive_hive = "2.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account_gen2 {
    storage_resource_id          = azurerm_storage_account.gen2test.id
    filesystem_id                = azurerm_storage_data_lake_gen2_filesystem.gen2test.id
    managed_identity_resource_id = azurerm_user_assigned_identity.test.id
    is_default                   = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D13_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_D14_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 2
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, template, rInt)
}

func testAccAzureRMHDInsightInteractiveQueryCluster_requiresImport(rInt int, rString string, location string) string {
	template := testAccAzureRMHDInsightInteractiveQueryCluster_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_interactive_query_cluster" "import" {
  name                = "${azurerm_hdinsight_interactive_query_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_interactive_query_cluster.test.resource_group_name}"
  location            = "${azurerm_hdinsight_interactive_query_cluster.test.location}"
  cluster_version     = "${azurerm_hdinsight_interactive_query_cluster.test.cluster_version}"
  tier                = "${azurerm_hdinsight_interactive_query_cluster.test.tier}"
  component_version   = "${azurerm_hdinsight_interactive_query_cluster.test.component_version}"
  gateway             = "${azurerm_hdinsight_interactive_query_cluster.test.gateway}"
  storage_account     = "${azurerm_hdinsight_interactive_query_cluster.test.storage_account}"
  roles               = "${azurerm_hdinsight_interactive_query_cluster.test.roles}"
}
`, template)
}

func testAccAzureRMHDInsightInteractiveQueryCluster_sshKeys(rInt int, rString string, location string) string {
	template := testAccAzureRMHDInsightInteractiveQueryCluster_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

variable "ssh_key" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
}

resource "azurerm_hdinsight_interactive_query_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"

  component_version {
    interactive_hive = "2.1"
  }

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
      vm_size  = "Standard_D13_V2"
      username = "acctestusrvm"
      ssh_keys = ["${var.ssh_key}"]
    }

    worker_node {
      vm_size               = "Standard_D14_V2"
      username              = "acctestusrvm"
      ssh_keys              = ["${var.ssh_key}"]
      target_instance_count = 3
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      ssh_keys = ["${var.ssh_key}"]
    }
  }
}
`, template, rInt)
}

func testAccAzureRMHDInsightInteractiveQueryCluster_updated(rInt int, rString string, location string) string {
	template := testAccAzureRMHDInsightInteractiveQueryCluster_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_interactive_query_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"

  component_version {
    interactive_hive = "2.1"
  }

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
      vm_size  = "Standard_D13_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_D14_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 5
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }

  tags = {
    Hello = "World"
  }
}
`, template, rInt)
}

func testAccAzureRMHDInsightInteractiveQueryCluster_virtualNetwork(rInt int, rString string, location string) string {
	template := testAccAzureRMHDInsightInteractiveQueryCluster_template(rInt, rString, location)
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

resource "azurerm_hdinsight_interactive_query_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"

  component_version {
    interactive_hive = "2.1"
  }

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
      vm_size            = "Standard_D13_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = "${azurerm_subnet.test.id}"
      virtual_network_id = "${azurerm_virtual_network.test.id}"
    }

    worker_node {
      vm_size               = "Standard_D14_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
      subnet_id             = "${azurerm_subnet.test.id}"
      virtual_network_id    = "${azurerm_virtual_network.test.id}"
    }

    zookeeper_node {
      vm_size            = "Standard_A4_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = "${azurerm_subnet.test.id}"
      virtual_network_id = "${azurerm_virtual_network.test.id}"
    }
  }
}
`, template, rInt, rInt, rInt)
}

func testAccAzureRMHDInsightInteractiveQueryCluster_complete(rInt int, rString string, location string) string {
	template := testAccAzureRMHDInsightInteractiveQueryCluster_template(rInt, rString, location)
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

resource "azurerm_hdinsight_interactive_query_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "3.6"
  tier                = "Standard"

  component_version {
    interactive_hive = "2.1"
  }

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
      vm_size            = "Standard_D13_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = "${azurerm_subnet.test.id}"
      virtual_network_id = "${azurerm_virtual_network.test.id}"
    }

    worker_node {
      vm_size               = "Standard_D14_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
      subnet_id             = "${azurerm_subnet.test.id}"
      virtual_network_id    = "${azurerm_virtual_network.test.id}"
    }

    zookeeper_node {
      vm_size            = "Standard_A4_V2"
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
`, template, rInt, rInt, rInt)
}

func testAccAzureRMHDInsightInteractiveQueryCluster_template(rInt int, rString string, location string) string {
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
`, rInt, location, rString)
}

func testAccAzureRMHDInsightInteractiveQueryCluster_gen2template(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "gen2test" {
  name                     = "accgen2test%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}

resource "azurerm_storage_data_lake_gen2_filesystem" "gen2test" {
  name               = "acctest"
  storage_account_id = azurerm_storage_account.gen2test.id
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  name = "test-identity"
}

data "azurerm_subscription" "primary" {}

resource "azurerm_role_assignment" "test" {
  scope                = "${data.azurerm_subscription.primary.id}"
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = "${azurerm_user_assigned_identity.test.principal_id}"
}
`, rInt, location, rString)
}
