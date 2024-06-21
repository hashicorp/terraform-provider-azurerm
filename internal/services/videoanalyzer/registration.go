// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package videoanalyzer

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/video-analyzer"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Video Analyzer"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Video Analyzer",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	if !features.FourPointOhBeta() {
		return map[string]*pluginsdk.Resource{
			"azurerm_video_analyzer":             resourceVideoAnalyzer(),
			"azurerm_video_analyzer_edge_module": resourceVideoAnalyzerEdgeModule(),
		}
	}

	return map[string]*pluginsdk.Resource{}
}
