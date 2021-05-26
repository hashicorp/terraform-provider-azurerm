package hdinsight

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "HDInsight"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"HDInsight",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_hdinsight_cluster": dataSourceHDInsightSparkCluster(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_hdinsight_hadoop_cluster":            resourceHDInsightHadoopCluster(),
		"azurerm_hdinsight_hbase_cluster":             resourceHDInsightHBaseCluster(),
		"azurerm_hdinsight_interactive_query_cluster": resourceHDInsightInteractiveQueryCluster(),
		"azurerm_hdinsight_kafka_cluster":             resourceHDInsightKafkaCluster(),
		"azurerm_hdinsight_ml_services_cluster":       resourceHDInsightMLServicesCluster(),
		"azurerm_hdinsight_rserver_cluster":           resourceHDInsightRServerCluster(),
		"azurerm_hdinsight_spark_cluster":             resourceHDInsightSparkCluster(),
		"azurerm_hdinsight_storm_cluster":             resourceHDInsightStormCluster(),
	}
}
