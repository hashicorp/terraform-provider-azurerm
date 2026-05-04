// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
		"azurerm_machine_learning_workspace": dataSourceMachineLearningWorkspace(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_machine_learning_compute_cluster":   resourceComputeCluster(),
		"azurerm_machine_learning_compute_instance":  resourceComputeInstance(),
		"azurerm_machine_learning_inference_cluster": resourceAksInferenceCluster(),
		"azurerm_machine_learning_synapse_spark":     resourceSynapseSpark(),
		"azurerm_machine_learning_workspace":         resourceMachineLearningWorkspace(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AIFoundry{},
		AIFoundryProject{},
		MachineLearningDataStoreBlobStorage{},
		MachineLearningDataStoreDataLakeGen2{},
		MachineLearningDataStoreFileShare{},
		WorkspaceNetworkOutboundRuleFqdn{},
		WorkspaceNetworkOutboundRulePrivateEndpoint{},
		WorkspaceNetworkOutboundRuleServiceTag{},
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
