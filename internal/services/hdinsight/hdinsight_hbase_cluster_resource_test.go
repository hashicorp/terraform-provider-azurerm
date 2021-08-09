package hdinsight_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HDInsightHBaseClusterResource struct {
}

func TestAccHDInsightHBaseCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHBaseCluster_gen2basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gen2basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHBaseCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccHDInsightHBaseCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHBaseCluster_sshKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sshKeys(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("storage_account",
			"roles.0.head_node.0.ssh_keys",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.ssh_keys",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.ssh_keys",
			"roles.0.zookeeper_node.0.vm_size"),
	})
}

func TestAccHDInsightHBaseCluster_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetwork(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHBaseCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHBaseCluster_tls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tls(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHBaseCluster_allMetastores(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allMetastores(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account",
			"metastores.0.hive.0.password",
			"metastores.0.oozie.0.password",
			"metastores.0.ambari.0.password"),
	})
}

func TestAccHDInsightHBaseCluster_hiveMetastore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hiveMetastore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
	})
}

func TestAccHDInsightHBaseCluster_updateMetastore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hiveMetastore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account",
			"metastores.0.hive.0.password",
			"metastores.0.oozie.0.password",
			"metastores.0.ambari.0.password"),
		{
			Config: r.allMetastores(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account",
			"metastores.0.hive.0.password",
			"metastores.0.oozie.0.password",
			"metastores.0.ambari.0.password"),
	})
}

func TestAccHDInsightHBaseCluster_monitor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.monitor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHBaseCluster_updateMonitor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		// No monitor
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
		// Add monitor
		{
			Config: r.monitor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
		// Change Log Analytics Workspace for the monitor
		{
			PreConfig: func() {
				data.RandomString += "new"
			},
			Config: r.monitor(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
		// Remove monitor
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccAzureRMHDInsightHBaseCluster_autoscale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoscale_schedule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHBaseCluster_securityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.securityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account"),
	})
}

func (t HDInsightHBaseClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resourceGroup := id.ResourceGroup
	name := id.Name

	resp, err := clients.HDInsight.ClustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading HDInsight HBase Cluster (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r HDInsightHBaseClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hbase = "2.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
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
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) gen2basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hbase = "2.1"
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
  }
}
`, r.gen2template(data), data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hbase_cluster" "import" {
  name                = azurerm_hdinsight_hbase_cluster.test.name
  resource_group_name = azurerm_hdinsight_hbase_cluster.test.resource_group_name
  location            = azurerm_hdinsight_hbase_cluster.test.location
  cluster_version     = azurerm_hdinsight_hbase_cluster.test.cluster_version
  tier                = azurerm_hdinsight_hbase_cluster.test.tier
  dynamic "component_version" {
    for_each = azurerm_hdinsight_hbase_cluster.test.component_version
    content {
      hbase = component_version.value.hbase
    }
  }
  dynamic "gateway" {
    for_each = azurerm_hdinsight_hbase_cluster.test.gateway
    content {
      enabled  = gateway.value.enabled
      password = gateway.value.password
      username = gateway.value.username
    }
  }
  dynamic "storage_account" {
    for_each = azurerm_hdinsight_hbase_cluster.test.storage_account
    content {
      is_default           = storage_account.value.is_default
      storage_account_key  = storage_account.value.storage_account_key
      storage_container_id = storage_account.value.storage_container_id
    }
  }
  dynamic "roles" {
    for_each = azurerm_hdinsight_hbase_cluster.test.roles
    content {
      dynamic "head_node" {
        for_each = lookup(roles.value, "head_node", [])
        content {
          password           = lookup(head_node.value, "password", null)
          subnet_id          = lookup(head_node.value, "subnet_id", null)
          username           = head_node.value.username
          virtual_network_id = lookup(head_node.value, "virtual_network_id", null)
          vm_size            = head_node.value.vm_size
        }
      }

      dynamic "worker_node" {
        for_each = lookup(roles.value, "worker_node", [])
        content {
          password              = lookup(worker_node.value, "password", null)
          subnet_id             = lookup(worker_node.value, "subnet_id", null)
          target_instance_count = worker_node.value.target_instance_count
          username              = worker_node.value.username
          virtual_network_id    = lookup(worker_node.value, "virtual_network_id", null)
          vm_size               = worker_node.value.vm_size
        }
      }

      dynamic "zookeeper_node" {
        for_each = lookup(roles.value, "zookeeper_node", [])
        content {
          password           = lookup(zookeeper_node.value, "password", null)
          subnet_id          = lookup(zookeeper_node.value, "subnet_id", null)
          username           = zookeeper_node.value.username
          virtual_network_id = lookup(zookeeper_node.value, "virtual_network_id", null)
          vm_size            = zookeeper_node.value.vm_size
        }
      }
    }
  }
}
`, r.basic(data))
}

func (r HDInsightHBaseClusterResource) sshKeys(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

variable "ssh_key" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
}

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hbase = "2.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
    is_default           = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      ssh_keys = [var.ssh_key]
    }

    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      ssh_keys              = [var.ssh_key]
      target_instance_count = 3
    }

    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      ssh_keys = [var.ssh_key]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hbase = "2.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
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
  }

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) virtualNetwork(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hbase = "2.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
    is_default           = true
  }

  roles {
    head_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
      subnet_id             = azurerm_subnet.test.id
      virtual_network_id    = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hbase = "2.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
    is_default           = true
  }

  roles {
    head_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
      subnet_id             = azurerm_subnet.test.id
      virtual_network_id    = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_D3_V2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (HDInsightHBaseClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (HDInsightHBaseClusterResource) gen2template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "gen2test" {
  depends_on = [azurerm_role_assignment.test]

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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r HDInsightHBaseClusterResource) tls(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    hbase = "2.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
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
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) allMetastores(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_sql_server" "test" {
  name                         = "acctestsql-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "sql_admin"
  administrator_login_password = "TerrAform123!"
  version                      = "12.0"
}
resource "azurerm_sql_database" "hive" {
  name                             = "hive"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  server_name                      = azurerm_sql_server.test.name
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  create_mode                      = "Default"
  requested_service_objective_name = "GP_Gen5_2"
}
resource "azurerm_sql_database" "oozie" {
  name                             = "oozie"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  server_name                      = azurerm_sql_server.test.name
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  create_mode                      = "Default"
  requested_service_objective_name = "GP_Gen5_2"
}
resource "azurerm_sql_database" "ambari" {
  name                             = "ambari"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  server_name                      = azurerm_sql_server.test.name
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  create_mode                      = "Default"
  requested_service_objective_name = "GP_Gen5_2"
}
resource "azurerm_sql_firewall_rule" "AzureServices" {
  name                = "allow-azure-services"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}
resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"
  component_version {
    hbase = "2.1"
  }
  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }
  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
    is_default           = true
  }
  roles {
    head_node {
      vm_size  = "Standard_D3_v2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
    worker_node {
      vm_size               = "Standard_D4_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 2
    }
    zookeeper_node {
      vm_size  = "Standard_D3_v2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
  metastores {
    hive {
      server        = azurerm_sql_server.test.fully_qualified_domain_name
      database_name = azurerm_sql_database.hive.name
      username      = azurerm_sql_server.test.administrator_login
      password      = azurerm_sql_server.test.administrator_login_password
    }
    oozie {
      server        = azurerm_sql_server.test.fully_qualified_domain_name
      database_name = azurerm_sql_database.oozie.name
      username      = azurerm_sql_server.test.administrator_login
      password      = azurerm_sql_server.test.administrator_login_password
    }
    ambari {
      server        = azurerm_sql_server.test.fully_qualified_domain_name
      database_name = azurerm_sql_database.ambari.name
      username      = azurerm_sql_server.test.administrator_login
      password      = azurerm_sql_server.test.administrator_login_password
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) hiveMetastore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_sql_server" "test" {
  name                         = "acctestsql-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "sql_admin"
  administrator_login_password = "TerrAform123!"
  version                      = "12.0"
}
resource "azurerm_sql_database" "hive" {
  name                             = "hive"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  server_name                      = azurerm_sql_server.test.name
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  create_mode                      = "Default"
  requested_service_objective_name = "GP_Gen5_2"
}
resource "azurerm_sql_firewall_rule" "AzureServices" {
  name                = "allow-azure-services"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}
resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"
  component_version {
    hbase = "2.1"
  }
  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }
  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
    is_default           = true
  }
  roles {
    head_node {
      vm_size  = "Standard_D3_v2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
    worker_node {
      vm_size               = "Standard_D4_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 2
    }
    zookeeper_node {
      vm_size  = "Standard_D3_v2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
  metastores {
    hive {
      server        = azurerm_sql_server.test.fully_qualified_domain_name
      database_name = azurerm_sql_database.hive.name
      username      = azurerm_sql_server.test.administrator_login
      password      = azurerm_sql_server.test.administrator_login_password
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) monitor(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%s-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hbase = "2.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
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
  }

  monitor {
    log_analytics_workspace_id = azurerm_log_analytics_workspace.test.workspace_id
    primary_key                = azurerm_log_analytics_workspace.test.primary_shared_key
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) autoscale_schedule(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hbase = "2.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
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
      autoscale {
        recurrence {
          timezone = "Pacific Standard Time"
          schedule {
            days                  = ["Monday"]
            time                  = "10:00"
            target_instance_count = 5
          }
          schedule {
            days                  = ["Saturday", "Sunday"]
            time                  = "10:00"
            target_instance_count = 3
          }
        }
      }
    }

    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) securityProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRGhdi-%d"
  location = "%s"

  tags = {
    StorageType = "Standard_LRS"
    type        = "test"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc-%d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.10.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSubnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = [cidrsubnet(azurerm_virtual_network.test.address_space.0, 8, 0)]
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestNSG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "AllowSyncWithAzureAD"
    priority                   = 101
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowRD"
    priority                   = 201
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3389"
    source_address_prefix      = "CorpNetSaw"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowPSRemoting"
    priority                   = 301
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "5986"
    source_address_prefix      = "AzureActiveDirectoryDomainServices"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowLDAPS"
    priority                   = 401
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "636"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource azurerm_subnet_network_security_group_association "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azuread_group" "test" {
  display_name = "AAD DC Administrators %s"
  description  = "Test for delegating group to administer Azure AD Domain Services"
}

data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test" {
  user_principal_name = "acctestAADDSAdminUser-%s@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestAADDSAdminUser-%s"
  password            = "TerrAform321!"
}

resource "azuread_group_member" "test" {
  group_object_id  = azuread_group.test.object_id
  member_object_id = azuread_user.test.object_id
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

data "azurerm_subscription" "primary" {}

resource "azurerm_role_assignment" "test" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "HDInsight Domain Services Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_active_directory_domain_service" "test" {
  name                = "acctestAADDS-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  domain_name           = "never.gonna.shut.you.down"
  sku                   = "Enterprise"
  filtered_sync_enabled = false

  initial_replica_set {
    subnet_id = azurerm_subnet.test.id
  }

  notifications {
    additional_recipients = ["notifyA@example.net", "notifyB@example.org"]
    notify_dc_admins      = true
    notify_global_admins  = true
  }

  secure_ldap {
    enabled                  = true
    external_access_enabled  = true
    pfx_certificate          = "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"
    pfx_certificate_password = "qwer5678"
  }

  security {
    ntlm_v1_enabled         = true
    sync_kerberos_passwords = true
    sync_ntlm_passwords     = true
    sync_on_prem_passwords  = true
    tls_v1_enabled          = true
  }

  tags = {
    Environment = "test"
  }

  depends_on = [
    azuread_group_member.test,
    azurerm_role_assignment.test,
    azurerm_subnet_network_security_group_association.test,
  ]
}

resource "azurerm_virtual_network_dns_servers" "test" {
  virtual_network_id = azurerm_virtual_network.test.id
  dns_servers        = azurerm_active_directory_domain_service.test.initial_replica_set.0.domain_controller_ip_addresses
}

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdihbase-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Premium"

  component_version {
    hbase = "2.1"
  }

  gateway {
    username = "sshuser"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.test.id
    storage_account_key  = azurerm_storage_account.test.primary_access_key
    is_default           = true
  }

  roles {
    head_node {
      vm_size            = "Standard_E4_V3"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_D12_V2"
      username              = "sshuser"
      password              = "TerrAform123!"
      target_instance_count = 2
      subnet_id             = azurerm_subnet.test.id
      virtual_network_id    = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_A2_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }

  security_profile {
    aadds_resource_id       = "${replace(azurerm_active_directory_domain_service.test.id, "/initialReplicaSetId/${azurerm_active_directory_domain_service.test.deployment_id}", "")}"
    domain_name             = azurerm_active_directory_domain_service.test.domain_name
    domain_username         = azuread_user.test.user_principal_name
    domain_user_password    = azuread_user.test.password
    ldaps_urls              = ["ldaps://${azurerm_active_directory_domain_service.test.domain_name}:636"]
    msi_resource_id         = azurerm_user_assigned_identity.test.id
    cluster_users_group_dns = [azuread_group.test.display_name]
  }

  depends_on = [
    azurerm_virtual_network_dns_servers.test,
  ]

  lifecycle {
    ignore_changes = [
      security_profile.0.domain_user_password,
      gateway.0.password
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomInteger)
}
