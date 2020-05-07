package machinelearning

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_machine_learning_workspace": dataSourceArmMachineLearningWorkspace(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_machine_learning_workspace": resourceArmMachineLearningWorkspace(),
	}
}
