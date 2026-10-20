// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/ai_foundry"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/ai_foundry_project"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_compute_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_compute_instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_datastore_blobstorage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_datastore_datalake_gen2"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_datastore_fileshare"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_inference_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_synapse_spark"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_workspace_network_outbound_rule_fqdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_workspace_network_outbound_rule_private_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/resources/machine_learning_workspace_network_outbound_rule_service_tag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/machine-learning"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Machine Learning"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Machine Learning",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_machine_learning_workspace": machine_learning_workspace.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_machine_learning_compute_cluster":   machine_learning_compute_cluster.RegisterResource(),
		"azurerm_machine_learning_compute_instance":  machine_learning_compute_instance.RegisterResource(),
		"azurerm_machine_learning_inference_cluster": machine_learning_inference_cluster.RegisterResource(),
		"azurerm_machine_learning_synapse_spark":     machine_learning_synapse_spark.RegisterResource(),
		"azurerm_machine_learning_workspace":         machine_learning_workspace.RegisterResource(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ai_foundry.Resource{},
		ai_foundry_project.Resource{},
		machine_learning_datastore_blobstorage.Resource{},
		machine_learning_datastore_datalake_gen2.Resource{},
		machine_learning_datastore_fileshare.Resource{},
		machine_learning_workspace_network_outbound_rule_fqdn.Resource{},
		machine_learning_workspace_network_outbound_rule_private_endpoint.Resource{},
		machine_learning_workspace_network_outbound_rule_service_tag.Resource{},
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
