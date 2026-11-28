// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_function_javascript_uda"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_function_javascript_udf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_job"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_job_schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_job_storage_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_managed_private_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_blob"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_eventhub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_function"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_mssql"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_powerbi"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_servicebus_queue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_servicebus_topic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_synapse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_table"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_reference_input_blob"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_reference_input_mssql"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_stream_input_blob"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_stream_input_eventhub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_stream_input_eventhub_v2"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_stream_input_iothub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/stream-analytics"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		stream_analytics_cluster.Resource{},
		stream_analytics_job_schedule.Resource{},
		stream_analytics_job_storage_account.Resource{},
		stream_analytics_managed_private_endpoint.Resource{},
		stream_analytics_output_cosmosdb.Resource{},
		stream_analytics_output_function.Resource{},
		stream_analytics_output_powerbi.Resource{},
		stream_analytics_output_table.Resource{},
		stream_analytics_stream_input_eventhub_v2.Resource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Stream Analytics"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Stream Analytics",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_stream_analytics_job": stream_analytics_job.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_stream_analytics_function_javascript_uda": stream_analytics_function_javascript_uda.RegisterResource(),
		"azurerm_stream_analytics_function_javascript_udf": stream_analytics_function_javascript_udf.RegisterResource(),
		"azurerm_stream_analytics_job":                     stream_analytics_job.RegisterResource(),
		"azurerm_stream_analytics_output_blob":             stream_analytics_output_blob.RegisterResource(),
		"azurerm_stream_analytics_output_eventhub":         stream_analytics_output_eventhub.RegisterResource(),
		"azurerm_stream_analytics_output_mssql":            stream_analytics_output_mssql.RegisterResource(),
		"azurerm_stream_analytics_output_servicebus_queue": stream_analytics_output_servicebus_queue.RegisterResource(),
		"azurerm_stream_analytics_output_servicebus_topic": stream_analytics_output_servicebus_topic.RegisterResource(),
		"azurerm_stream_analytics_output_synapse":          stream_analytics_output_synapse.RegisterResource(),
		"azurerm_stream_analytics_reference_input_blob":    stream_analytics_reference_input_blob.RegisterResource(),
		"azurerm_stream_analytics_reference_input_mssql":   stream_analytics_reference_input_mssql.RegisterResource(),
		"azurerm_stream_analytics_stream_input_blob":       stream_analytics_stream_input_blob.RegisterResource(),
		"azurerm_stream_analytics_stream_input_eventhub":   stream_analytics_stream_input_eventhub.RegisterResource(),
		"azurerm_stream_analytics_stream_input_iothub":     stream_analytics_stream_input_iothub.RegisterResource(),
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
