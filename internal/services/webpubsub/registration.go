package webpubsub

import "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

type Registration struct{}

func (r Registration) Name() string {
	return "Web Publishing Subscription"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Web Publishing Subscription",
	}
}
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_web_pubsub":     resourceWebPubSub(),
		"azurerm_web_pubsub_hub": resourceWebPubsubHub(),
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_web_pubsub": dataSourceWebPubsub(),
	}
}
