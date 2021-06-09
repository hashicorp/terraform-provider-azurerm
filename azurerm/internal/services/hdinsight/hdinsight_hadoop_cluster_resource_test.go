package hdinsight_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type HDInsightHadoopClusterResource struct {
}

func TestAccHDInsightHadoopCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}

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

func TestAccHDInsightHadoopCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_sshKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_edgeNodeBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.edgeNodeBasic(data, 2, "Standard_D3_V2"),
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
			"roles.0.edge_node.0.password",
			"roles.0.edge_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHadoopCluster_addEdgeNodeBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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
			Config: r.edgeNodeBasic(data, 1, "Standard_D3_V2"),
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
			"roles.0.edge_node.0.password",
			"roles.0.edge_node.0.vm_size",
			"storage_account"),
		{
			Config: r.edgeNodeBasic(data, 3, "Standard_D4_V2"),
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
			"roles.0.edge_node.0.password",
			"roles.0.edge_node.0.vm_size",
			"storage_account"),
	})
}

func TestAccHDInsightHadoopCluster_gen2storage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gen2storage(data),
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

func TestAccHDInsightHadoopCluster_gen2AndBlobStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gen2AndBlobStorage(data),
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

func TestAccHDInsightHadoopCluster_privateLink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateLink(data),
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

func TestAccHDInsightHadoopCluster_tls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_allMetastores(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_hiveMetastore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_updateMetastore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_monitor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccHDInsightHadoopCluster_updateGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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
			Config: r.updateGateway(data),
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

func TestAccHDInsightHadoopCluster_updateMonitor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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

func TestAccAzureRMHDInsightHadoopCluster_autoscale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}
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
			Config: r.autoscale_capacity(data),
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

func (t HDInsightHadoopClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resourceGroup := id.ResourceGroup
	name := id.Name

	resp, err := clients.HDInsight.ClustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading HDInsight Hadoop Cluster (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r HDInsightHadoopClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hadoop = "3.1"
  }

  gateway {
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
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hadoop_cluster" "import" {
  name                = azurerm_hdinsight_hadoop_cluster.test.name
  resource_group_name = azurerm_hdinsight_hadoop_cluster.test.resource_group_name
  location            = azurerm_hdinsight_hadoop_cluster.test.location
  cluster_version     = azurerm_hdinsight_hadoop_cluster.test.cluster_version
  tier                = azurerm_hdinsight_hadoop_cluster.test.tier
  dynamic "component_version" {
    for_each = azurerm_hdinsight_hadoop_cluster.test.component_version
    content {
      hadoop = component_version.value.hadoop
    }
  }
  dynamic "gateway" {
    for_each = azurerm_hdinsight_hadoop_cluster.test.gateway
    content {
      enabled  = gateway.value.enabled
      password = gateway.value.password
      username = gateway.value.username
    }
  }
  dynamic "storage_account" {
    for_each = azurerm_hdinsight_hadoop_cluster.test.storage_account
    content {
      is_default           = storage_account.value.is_default
      storage_account_key  = storage_account.value.storage_account_key
      storage_container_id = storage_account.value.storage_container_id
    }
  }
  dynamic "roles" {
    for_each = azurerm_hdinsight_hadoop_cluster.test.roles
    content {
      dynamic "edge_node" {
        for_each = lookup(roles.value, "edge_node", [])
        content {
          target_instance_count = edge_node.value.target_instance_count
          vm_size               = edge_node.value.vm_size

          dynamic "install_script_action" {
            for_each = lookup(edge_node.value, "install_script_action", [])
            content {
              name = install_script_action.value.name
              uri  = install_script_action.value.uri
            }
          }
        }
      }

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

func (r HDInsightHadoopClusterResource) sshKeys(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

variable "ssh_key" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
}

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hadoop = "3.1"
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
      ssh_keys = [var.ssh_key]
    }

    worker_node {
      vm_size               = "Standard_D4_v2"
      username              = "acctestusrvm"
      ssh_keys              = [var.ssh_key]
      target_instance_count = 3
    }

    zookeeper_node {
      vm_size  = "Standard_D3_v2"
      username = "acctestusrvm"
      ssh_keys = [var.ssh_key]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hadoop = "3.1"
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
      vm_size               = "Standard_D4_v2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 5
    }

    zookeeper_node {
      vm_size  = "Standard_D3_v2"
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

func (r HDInsightHadoopClusterResource) virtualNetwork(data acceptance.TestData) string {
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

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hadoop = "3.1"
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
      vm_size            = "Standard_D3_v2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_D4_v2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
      subnet_id             = azurerm_subnet.test.id
      virtual_network_id    = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_D3_v2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) complete(data acceptance.TestData) string {
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

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hadoop = "3.1"
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
      vm_size            = "Standard_D3_v2"
      username           = "acctestusrvm"
      password           = "AccTestvdSC4daf986!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_D4_v2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
      subnet_id             = azurerm_subnet.test.id
      virtual_network_id    = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_D3_v2"
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

func (r HDInsightHadoopClusterResource) edgeNodeBasic(data acceptance.TestData, numEdgeNodes int, instanceType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hadoop = "3.1"
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

    edge_node {
      target_instance_count = %d
      vm_size               = "%s"
      install_script_action {
        name = "script1"
        uri  = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-hdinsight-linux-with-edge-node/scripts/EmptyNodeSetup.sh"
      }
    }
  }
}
`, r.template(data), data.RandomInteger, numEdgeNodes, instanceType)
}

func (r HDInsightHadoopClusterResource) gen2storage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_hdinsight_hadoop_cluster" "test" {
  depends_on = [azurerm_role_assignment.test]

  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "4.0"
  tier                = "Standard"
  component_version {
    hadoop = "3.1"
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
}
`, r.gen2template(data), data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) gen2AndBlobStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  depends_on = [azurerm_role_assignment.test]

  name                = "acctesthdi-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "4.0"
  tier                = "Standard"
  component_version {
    hadoop = "3.1"
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
  storage_account {
    storage_container_id = "${azurerm_storage_container.test.id}"
    storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
    is_default           = false
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
}
`, r.gen2template(data), data.RandomString, data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) privateLink(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["172.16.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.16.11.0/26"]

  enforce_private_link_service_network_policies = true
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_nat_gateway" "test" {
  name                    = "acctestnat%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  zones                   = ["1"]
}

resource "azurerm_nat_gateway_public_ip_association" "test" {
  nat_gateway_id       = azurerm_nat_gateway.test.id
  public_ip_address_id = azurerm_public_ip.test.id
}

resource "azurerm_subnet_nat_gateway_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  nat_gateway_id = azurerm_nat_gateway.test.id
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  depends_on = [azurerm_role_assignment.test, azurerm_nat_gateway.test, azurerm_subnet_network_security_group_association.test]

  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hadoop = "3.1"
  }

  network {
    connection_direction = "Outbound"
    private_link_enabled = true
  }

  gateway {
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

      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_D4_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3

      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"

      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }
}

%s
`, r.gen2template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, r.nsgTemplate(data))
}

func (HDInsightHadoopClusterResource) nsgTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule = [
    {
      access                                     = "Allow"
      description                                = "Rule can be deleted but do not change source ips."
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = "443"
      destination_port_ranges                    = []
      direction                                  = "Inbound"
      name                                       = "Rule-101"
      priority                                   = 101
      protocol                                   = "Tcp"
      source_address_prefix                      = "VirtualNetwork"
      source_address_prefixes                    = []
      source_application_security_group_ids      = []
      source_port_range                          = "*"
      source_port_ranges                         = []
    },
    {
      access                                     = "Allow"
      description                                = "Rule can be deleted but do not change source ips."
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = "*"
      destination_port_ranges                    = []
      direction                                  = "Inbound"
      name                                       = "Rule-103"
      priority                                   = 103
      protocol                                   = "*"
      source_address_prefix                      = "CorpNetPublic"
      source_address_prefixes                    = []
      source_application_security_group_ids      = []
      source_port_range                          = "*"
      source_port_ranges                         = []
    },
    {
      access                                     = "Allow"
      description                                = "Rule can be deleted but do not change source ips."
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = "*"
      destination_port_ranges                    = []
      direction                                  = "Inbound"
      name                                       = "Rule-104"
      priority                                   = 104
      protocol                                   = "*"
      source_address_prefix                      = "CorpNetSaw"
      source_address_prefixes                    = []
      source_application_security_group_ids      = []
      source_port_range                          = "*"
      source_port_ranges                         = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "111",
        "11211",
        "123",
        "13",
        "17",
        "19",
        "1900",
        "512",
        "514",
        "53",
        "5353",
        "593",
        "69",
        "873",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-108"
      priority                              = 108
      protocol                              = "*"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "119",
        "137",
        "138",
        "139",
        "161",
        "162",
        "2049",
        "2301",
        "2381",
        "3268",
        "389",
        "5800",
        "5900",
        "636",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-109"
      priority                              = 109
      protocol                              = "*"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "135",
        "23",
        "445",
        "5985",
        "5986",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-107"
      priority                              = 107
      protocol                              = "Tcp"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "1433",
        "1434",
        "16379",
        "26379",
        "27017",
        "3306",
        "4333",
        "5432",
        "6379",
        "7000",
        "7001",
        "7199",
        "9042",
        "9160",
        "9300",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-105"
      priority                              = 105
      protocol                              = "*"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
    {
      access                                     = "Deny"
      description                                = "DO NOT DELETE"
      destination_address_prefix                 = "*"
      destination_address_prefixes               = []
      destination_application_security_group_ids = []
      destination_port_range                     = ""
      destination_port_ranges = [
        "22",
        "3389",
      ]
      direction                             = "Inbound"
      name                                  = "Rule-106"
      priority                              = 106
      protocol                              = "Tcp"
      source_address_prefix                 = "Internet"
      source_address_prefixes               = []
      source_application_security_group_ids = []
      source_port_range                     = "*"
      source_port_ranges                    = []
    },
  ]
}

`, data.RandomInteger)
}

func (HDInsightHadoopClusterResource) template(data acceptance.TestData) string {
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

func (HDInsightHadoopClusterResource) gen2template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r HDInsightHadoopClusterResource) tls(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"
  tls_min_version     = "1.2"
  component_version {
    hadoop = "3.1"
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
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) allMetastores(data acceptance.TestData) string {
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
resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"
  component_version {
    hadoop = "3.1"
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

func (r HDInsightHadoopClusterResource) hiveMetastore(data acceptance.TestData) string {
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
resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"
  component_version {
    hadoop = "3.1"
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

func (r HDInsightHadoopClusterResource) monitor(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%s-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    hadoop = "3.1"
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

  monitor {
    log_analytics_workspace_id = azurerm_log_analytics_workspace.test.workspace_id
    primary_key                = azurerm_log_analytics_workspace.test.primary_shared_key
  }
}
`, r.template(data), data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) updateGateway(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"
  component_version {
    hadoop = "3.1"
  }
  gateway {
    username = "acctestusrgw"
    password = "TerrAformne3!"
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
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) autoscale_capacity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"
  component_version {
    hadoop = "3.1"
  }
  gateway {
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
      autoscale {
        capacity {
          min_instance_count = 2
          max_instance_count = 3
        }
      }
    }
    zookeeper_node {
      vm_size  = "Standard_D3_v2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) autoscale_schedule(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "4.0"
  tier                = "Standard"
  component_version {
    hadoop = "3.1"
  }
  gateway {
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
      vm_size  = "Standard_D3_v2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, r.template(data), data.RandomInteger)
}
