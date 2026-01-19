// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arckubernetes

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

// Name is the name of this Service
func (r Registration) Name() string {
	return "ArcKubernetes"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"ArcKubernetes",
	}
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/arc-kubernetes"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_arc_kubernetes_cluster": resourceArcKubernetesCluster(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ArcKubernetesClusterExtensionResource{},
		ArcKubernetesFluxConfigurationResource{},
		ArcKubernetesProvisionedClusterResource{},
	}
}
