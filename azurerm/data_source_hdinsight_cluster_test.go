package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMHDInsightCluster_hadoop(t *testing.T) {
	dataSourceName := "data.azurerm_hdinsight_cluster.test"
	rInt := tf.AccRandTimeInt()
	rString := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_hadoop(rInt, rString, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "hadoop"),
					resource.TestCheckResourceAttr(dataSourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(dataSourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_hbase(t *testing.T) {
	dataSourceName := "data.azurerm_hdinsight_cluster.test"
	rInt := tf.AccRandTimeInt()
	rString := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_hbase(rInt, rString, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "hbase"),
					resource.TestCheckResourceAttr(dataSourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(dataSourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_interactiveQuery(t *testing.T) {
	dataSourceName := "data.azurerm_hdinsight_cluster.test"
	rInt := tf.AccRandTimeInt()
	rString := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_interactiveQuery(rInt, rString, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "interactivehive"),
					resource.TestCheckResourceAttr(dataSourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(dataSourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_kafka(t *testing.T) {
	dataSourceName := "data.azurerm_hdinsight_cluster.test"
	rInt := tf.AccRandTimeInt()
	rString := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_kafka(rInt, rString, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "kafka"),
					resource.TestCheckResourceAttr(dataSourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(dataSourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_mlServices(t *testing.T) {
	dataSourceName := "data.azurerm_hdinsight_cluster.test"
	rInt := tf.AccRandTimeInt()
	rString := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_mlServices(rInt, rString, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "mlservices"),
					resource.TestCheckResourceAttr(dataSourceName, "tier", "standard"),
					resource.TestCheckResourceAttrSet(dataSourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_rserver(t *testing.T) {
	dataSourceName := "data.azurerm_hdinsight_cluster.test"
	rInt := tf.AccRandTimeInt()
	rString := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_rserver(rInt, rString, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "rserver"),
					resource.TestCheckResourceAttr(dataSourceName, "tier", "standard"),
					resource.TestCheckResourceAttrSet(dataSourceName, "edge_ssh_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_spark(t *testing.T) {
	dataSourceName := "data.azurerm_hdinsight_cluster.test"
	rInt := tf.AccRandTimeInt()
	rString := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_spark(rInt, rString, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "spark"),
					resource.TestCheckResourceAttr(dataSourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(dataSourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMHDInsightCluster_storm(t *testing.T) {
	dataSourceName := "data.azurerm_hdinsight_cluster.test"
	rInt := tf.AccRandTimeInt()
	rString := strings.ToLower(acctest.RandString(11))
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHDInsightCluster_storm(rInt, rString, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "kind", "storm"),
					resource.TestCheckResourceAttr(dataSourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(dataSourceName, "edge_ssh_endpoint", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "https_endpoint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ssh_endpoint"),
				),
			},
		},
	})
}

func testAccDataSourceHDInsightCluster_hadoop(rInt int, rString, location string) string {
	template := testAccAzureRMHDInsightHadoopCluster_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = "${azurerm_hdinsight_hadoop_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_hadoop_cluster.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceHDInsightCluster_hbase(rInt int, rString, location string) string {
	template := testAccAzureRMHDInsightHBaseCluster_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = "${azurerm_hdinsight_hbase_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_hbase_cluster.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceHDInsightCluster_interactiveQuery(rInt int, rString, location string) string {
	template := testAccAzureRMHDInsightInteractiveQueryCluster_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = "${azurerm_hdinsight_interactive_query_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_interactive_query_cluster.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceHDInsightCluster_kafka(rInt int, rString, location string) string {
	template := testAccAzureRMHDInsightKafkaCluster_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = "${azurerm_hdinsight_kafka_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_kafka_cluster.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceHDInsightCluster_mlServices(rInt int, rString, location string) string {
	template := testAccAzureRMHDInsightMLServicesCluster_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = "${azurerm_hdinsight_ml_services_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_ml_services_cluster.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceHDInsightCluster_rserver(rInt int, rString, location string) string {
	template := testAccAzureRMHDInsightRServerCluster_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = "${azurerm_hdinsight_rserver_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_rserver_cluster.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceHDInsightCluster_spark(rInt int, rString, location string) string {
	template := testAccAzureRMHDInsightSparkCluster_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = "${azurerm_hdinsight_spark_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_spark_cluster.test.resource_group_name}"
}
`, template)
}

func testAccDataSourceHDInsightCluster_storm(rInt int, rString, location string) string {
	template := testAccAzureRMHDInsightStormCluster_basic(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_hdinsight_cluster" "test" {
  name                = "${azurerm_hdinsight_storm_cluster.test.name}"
  resource_group_name = "${azurerm_hdinsight_storm_cluster.test.resource_group_name}"
}
`, template)
}
