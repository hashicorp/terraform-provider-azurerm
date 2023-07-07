// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/hdinsight"
}

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
		"azurerm_hdinsight_spark_cluster":             resourceHDInsightSparkCluster(),
	}
}
