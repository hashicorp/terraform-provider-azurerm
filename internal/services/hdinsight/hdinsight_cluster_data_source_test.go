// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type HDInsightClusterDataSourceResource struct{}

func TestAccDataSourceHDInsightCluster_hadoop(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.hadoop(data),
			Check: acceptance.ComposeTestCheckFunc(
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
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.hbase(data),
			Check: acceptance.ComposeTestCheckFunc(
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
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.interactiveQuery(data),
			Check: acceptance.ComposeTestCheckFunc(
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
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.kafka(data),
			Check: acceptance.ComposeTestCheckFunc(
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
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.kafkaWithRestProxy(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccDataSourceHDInsightCluster_spark(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	r := HDInsightClusterDataSourceResource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.spark(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("spark"),
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

func (HDInsightClusterDataSourceResource) spark(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_spark_cluster.test.name
  resource_group_name = azurerm_hdinsight_spark_cluster.test.resource_group_name
}
`, HDInsightSparkClusterResource{}.basic(data))
}
