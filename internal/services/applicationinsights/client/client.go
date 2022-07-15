package client

import (
	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2020-02-02/insights"
	workbookTemplate "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/applicationinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/azuresdkhacks"
)

type Client struct {
	AnalyticsItemsClient     *insights.AnalyticsItemsClient
	APIKeysClient            *insights.APIKeysClient
	ComponentsClient         *insights.ComponentsClient
	WebTestsClient           *azuresdkhacks.WebTestsClient
	BillingClient            *insights.ComponentCurrentBillingFeaturesClient
	SmartDetectionRuleClient *insights.ProactiveDetectionConfigurationsClient
	WorkbookTemplateClient   *workbookTemplate.ApplicationInsightsClient
}

func NewClient(o *common.ClientOptions) *Client {
	analyticsItemsClient := insights.NewAnalyticsItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&analyticsItemsClient.Client, o.ResourceManagerAuthorizer)

	apiKeysClient := insights.NewAPIKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiKeysClient.Client, o.ResourceManagerAuthorizer)

	componentsClient := insights.NewComponentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&componentsClient.Client, o.ResourceManagerAuthorizer)

	webTestsClient := insights.NewWebTestsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webTestsClient.Client, o.ResourceManagerAuthorizer)
	webTestsWorkaroundClient := azuresdkhacks.NewWebTestsClient(webTestsClient)

	billingClient := insights.NewComponentCurrentBillingFeaturesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&billingClient.Client, o.ResourceManagerAuthorizer)

	smartDetectionRuleClient := insights.NewProactiveDetectionConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&smartDetectionRuleClient.Client, o.ResourceManagerAuthorizer)

	workbookTemplateClient := workbookTemplate.NewApplicationInsightsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&workbookTemplateClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AnalyticsItemsClient:     &analyticsItemsClient,
		APIKeysClient:            &apiKeysClient,
		ComponentsClient:         &componentsClient,
		WebTestsClient:           &webTestsWorkaroundClient,
		BillingClient:            &billingClient,
		SmartDetectionRuleClient: &smartDetectionRuleClient,
		WorkbookTemplateClient:   &workbookTemplateClient,
	}
}
