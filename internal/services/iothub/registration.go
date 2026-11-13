// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_consumer_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_device_update_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_device_update_instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_dps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_dps_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_dps_shared_access_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_endpoint_cosmosdb_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_endpoint_eventhub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_endpoint_servicebus_queue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_endpoint_servicebus_topic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_endpoint_storage_container"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_enrichment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_fallback_route"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_file_upload"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_route"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/resources/iothub_shared_access_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/iot-hub"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "IoT Hub"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"IoT Hub",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_iothub":                          iothub.RegisterDataSource(),
		"azurerm_iothub_dps":                      iothub_dps.RegisterDataSource(),
		"azurerm_iothub_dps_shared_access_policy": iothub_dps_shared_access_policy.RegisterDataSource(),
		"azurerm_iothub_shared_access_policy":     iothub_shared_access_policy.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_iothub":                            iothub.RegisterResource(),
		"azurerm_iothub_certificate":                iothub_certificate.RegisterResource(),
		"azurerm_iothub_consumer_group":             iothub_consumer_group.RegisterResource(),
		"azurerm_iothub_dps":                        iothub_dps.RegisterResource(),
		"azurerm_iothub_dps_certificate":            iothub_dps_certificate.RegisterResource(),
		"azurerm_iothub_dps_shared_access_policy":   iothub_dps_shared_access_policy.RegisterResource(),
		"azurerm_iothub_endpoint_eventhub":          iothub_endpoint_eventhub.RegisterResource(),
		"azurerm_iothub_endpoint_servicebus_queue":  iothub_endpoint_servicebus_queue.RegisterResource(),
		"azurerm_iothub_endpoint_servicebus_topic":  iothub_endpoint_servicebus_topic.RegisterResource(),
		"azurerm_iothub_endpoint_storage_container": iothub_endpoint_storage_container.RegisterResource(),
		"azurerm_iothub_enrichment":                 iothub_enrichment.RegisterResource(),
		"azurerm_iothub_fallback_route":             iothub_fallback_route.RegisterResource(),
		"azurerm_iothub_route":                      iothub_route.RegisterResource(),
		"azurerm_iothub_shared_access_policy":       iothub_shared_access_policy.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		iothub_device_update_account.Resource{},
		iothub_device_update_instance.Resource{},
		iothub_endpoint_cosmosdb_account.Resource{},
		iothub_file_upload.Resource{},
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
