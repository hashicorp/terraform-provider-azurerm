package hdinsight_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMHDInsightCluster_hadoop(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_hadoop(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "hadoop"),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_hbase(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_hbase(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "hbase"),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_interactiveQuery(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_interactiveQuery(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "interactivehive"),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_kafka(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_kafka(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "kafka"),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_mlServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_mlServices(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "mlservices"),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "standard"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_rserver(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_rserver(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "rserver"),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "standard"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_spark(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_spark(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "spark"),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_storm(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_hdinsight_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_storm(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "storm"),
					resource.TestCheckResourceAttr(data.ResourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(data.ResourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func testAccDataSourceHDInsightCluster_hadoop(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightHadoopCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_hadoop_cluster.test.name
  resource_group_name = azurerm_hdinsight_hadoop_cluster.test.resource_group_name
}
`, template)
}

func testAccDataSourceHDInsightCluster_hbase(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightHBaseCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_hbase_cluster.test.name
  resource_group_name = azurerm_hdinsight_hbase_cluster.test.resource_group_name
}
`, template)
}

func testAccDataSourceHDInsightCluster_interactiveQuery(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightInteractiveQueryCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_interactive_query_cluster.test.name
  resource_group_name = azurerm_hdinsight_interactive_query_cluster.test.resource_group_name
}
`, template)
}

func testAccDataSourceHDInsightCluster_kafka(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightKafkaCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_kafka_cluster.test.name
  resource_group_name = azurerm_hdinsight_kafka_cluster.test.resource_group_name
}
`, template)
}

func testAccDataSourceHDInsightCluster_mlServices(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightMLServicesCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_ml_services_cluster.test.name
  resource_group_name = azurerm_hdinsight_ml_services_cluster.test.resource_group_name
}
`, template)
}

func testAccDataSourceHDInsightCluster_rserver(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightRServerCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_rserver_cluster.test.name
  resource_group_name = azurerm_hdinsight_rserver_cluster.test.resource_group_name
}
`, template)
}

func testAccDataSourceHDInsightCluster_spark(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightSparkCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_spark_cluster.test.name
  resource_group_name = azurerm_hdinsight_spark_cluster.test.resource_group_name
}
`, template)
}

func testAccDataSourceHDInsightCluster_storm(data acceptance.TestData) string {
	template := testAccAzureRMHDInsightStormCluster_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = azurerm_hdinsight_storm_cluster.test.name
  resource_group_name = azurerm_hdinsight_storm_cluster.test.resource_group_name
}
`, template)
}
