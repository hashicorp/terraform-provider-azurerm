package machinelearning

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

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
		"azurerm_machine_learning_workspace":         resourceMachineLearningWorkspace(),
		"azurerm_machine_learning_inference_cluster": resourceAksInferenceCluster(),
		"azurerm_machine_learning_compute_cluster":   resourceComputeCluster(),
	}
}
