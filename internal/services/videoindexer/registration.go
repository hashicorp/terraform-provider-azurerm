package videoindexer

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/video-indexer"
}

func (r Registration) WebsiteCategories() []string {
	return []string{"Video Indexer"}
}

func (r Registration) Name() string {
	return "VideoIndexer"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AccountResource{},
	}
}
