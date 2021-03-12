package hdinsight_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type HDInsightClusterDataSourceResource struct {
}

func TestAccDataSourceHDInsightCluster_hadoop(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.hadoop(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("hadoop"),
				check.That(data.ResourceName).Key("tier").HasValue("standard"),
				check.That(data.ResourceName).Key("edge_ssh_endpoint").HasValue(""),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceHDInsightCluster_hbase(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.hbase(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("hbase"),
				check.That(data.ResourceName).Key("tier").HasValue("standard"),
				check.That(data.ResourceName).Key("edge_ssh_endpoint").HasValue(""),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceHDInsightCluster_interactiveQuery(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.interactiveQuery(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("interactivehive"),
				check.That(data.ResourceName).Key("tier").HasValue("standard"),
				check.That(data.ResourceName).Key("edge_ssh_endpoint").HasValue(""),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceHDInsightCluster_kafka(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.kafka(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("kafka"),
				check.That(data.ResourceName).Key("tier").HasValue("standard"),
				check.That(data.ResourceName).Key("edge_ssh_endpoint").HasValue(""),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceHDInsightCluster_kafkaWithRestProxy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.kafkaWithRestProxy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("kafka"),
				check.That(data.ResourceName).Key("tier").HasValue("standard"),
				check.That(data.ResourceName).Key("edge_ssh_endpoint").HasValue(""),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
				check.That(data.ResourceName).Key("kafka_rest_proxy_endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceHDInsightCluster_mlServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.mlServices(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("mlservices"),
				check.That(data.ResourceName).Key("tier").HasValue("standard"),
				check.That(data.ResourceName).Key("edge_ssh_endpoint").Exists(),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceHDInsightCluster_rserver(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.rserver(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("rserver"),
				check.That(data.ResourceName).Key("tier").HasValue("standard"),
				check.That(data.ResourceName).Key("edge_ssh_endpoint").Exists(),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceHDInsightCluster_spark(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.spark(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("spark"),
				check.That(data.ResourceName).Key("tier").HasValue("standard"),
				check.That(data.ResourceName).Key("edge_ssh_endpoint").HasValue(""),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
	})
}

func TestAccDataSourceHDInsightCluster_storm(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.storm(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("storm"),
				check.That(data.ResourceName).Key("tier").HasValue("standard"),
				check.That(data.ResourceName).Key("edge_ssh_endpoint").HasValue(""),
				check.That(data.ResourceName).Key("https_endpoint").Exists(),
				check.That(data.ResourceName).Key("ssh_endpoint").Exists(),
			),
		},
	})
}

func (HDInsightClusterDataSourceResource) hadoop(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_hadoop_cluster.test.name
  resource_group_name = azurerm_hdinsight_hadoop_cluster.test.resource_group_name
}
`, HDInsightHadoopClusterResource{}.basic(data))
}

func (HDInsightClusterDataSourceResource) hbase(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_hbase_cluster.test.name
  resource_group_name = azurerm_hdinsight_hbase_cluster.test.resource_group_name
}
`, HDInsightHBaseClusterResource{}.basic(data))
}

func (HDInsightClusterDataSourceResource) interactiveQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_interactive_query_cluster.test.name
  resource_group_name = azurerm_hdinsight_interactive_query_cluster.test.resource_group_name
}
`, HDInsightInteractiveQueryClusterResource{}.basic(data))
}

func (HDInsightClusterDataSourceResource) kafka(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_kafka_cluster.test.name
  resource_group_name = azurerm_hdinsight_kafka_cluster.test.resource_group_name
}
`, HDInsightKafkaClusterResource{}.basic(data))
}

func (HDInsightClusterDataSourceResource) kafkaWithRestProxy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_kafka_cluster.test.name
  resource_group_name = azurerm_hdinsight_kafka_cluster.test.resource_group_name
}
`, HDInsightKafkaClusterResource{}.restProxy(data))
}

func (HDInsightClusterDataSourceResource) mlServices(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_ml_services_cluster.test.name
  resource_group_name = azurerm_hdinsight_ml_services_cluster.test.resource_group_name
}
`, HDInsightMLServicesClusterResource{}.basic(data))
}

func (HDInsightClusterDataSourceResource) rserver(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_rserver_cluster.test.name
  resource_group_name = azurerm_hdinsight_rserver_cluster.test.resource_group_name
}
`, HDInsightRServerClusterResource{}.basic(data))
}

func (HDInsightClusterDataSourceResource) spark(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_spark_cluster.test.name
  resource_group_name = azurerm_hdinsight_spark_cluster.test.resource_group_name
}
`, HDInsightSparkClusterResource{}.basic(data))
}

func (HDInsightClusterDataSourceResource) storm(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_storm_cluster.test.name
  resource_group_name = azurerm_hdinsight_storm_cluster.test.resource_group_name
}
`, HDInsightStormClusterResource{}.basic(data))
}
