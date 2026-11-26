// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/resources/hdinsight_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/resources/hdinsight_hadoop_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/resources/hdinsight_hbase_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/resources/hdinsight_interactive_query_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/resources/hdinsight_kafka_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/resources/hdinsight_spark_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

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
		"azurerm_hdinsight_cluster": hdinsight_cluster.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_hdinsight_hadoop_cluster":            hdinsight_hadoop_cluster.RegisterResource(),
		"azurerm_hdinsight_hbase_cluster":             hdinsight_hbase_cluster.RegisterResource(),
		"azurerm_hdinsight_interactive_query_cluster": hdinsight_interactive_query_cluster.RegisterResource(),
		"azurerm_hdinsight_kafka_cluster":             hdinsight_kafka_cluster.RegisterResource(),
		"azurerm_hdinsight_spark_cluster":             hdinsight_spark_cluster.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
