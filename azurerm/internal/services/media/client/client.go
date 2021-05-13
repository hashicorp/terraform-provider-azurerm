package client

import (
	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ServicesClient           *media.MediaservicesClient
	AssetsClient             *media.AssetsClient
	TransformsClient         *media.TransformsClient
	StreamingEndpointsClient *media.StreamingEndpointsClient
	JobsClient               *media.JobsClient
	StreamingLocatorsClient  *media.StreamingLocatorsClient
	ContentKeyPoliciesClient *media.ContentKeyPoliciesClient
	StreamingPoliciesClient  *media.StreamingPoliciesClient
	LiveEventsClient         *media.LiveEventsClient
	LiveOutputsClient        *media.LiveOutputsClient
	AssetFiltersClient       *media.AssetFiltersClient
}

func NewClient(o *common.ClientOptions) *Client {
	ServicesClient := media.NewMediaservicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServicesClient.Client, o.ResourceManagerAuthorizer)

	AssetsClient := media.NewAssetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AssetsClient.Client, o.ResourceManagerAuthorizer)

	TransformsClient := media.NewTransformsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TransformsClient.Client, o.ResourceManagerAuthorizer)

	StreamingEndpointsClient := media.NewStreamingEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&StreamingEndpointsClient.Client, o.ResourceManagerAuthorizer)

	JobsClient := media.NewJobsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&JobsClient.Client, o.ResourceManagerAuthorizer)

	StreamingLocatorsClient := media.NewStreamingLocatorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&StreamingLocatorsClient.Client, o.ResourceManagerAuthorizer)

	ContentKeyPoliciesClient := media.NewContentKeyPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ContentKeyPoliciesClient.Client, o.ResourceManagerAuthorizer)

	StreamingPoliciesClient := media.NewStreamingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&StreamingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	LiveEventsClient := media.NewLiveEventsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LiveEventsClient.Client, o.ResourceManagerAuthorizer)

	LiveOutputsClient := media.NewLiveOutputsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LiveOutputsClient.Client, o.ResourceManagerAuthorizer)

	AssetFiltersClient := media.NewAssetFiltersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AssetFiltersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServicesClient:           &ServicesClient,
		AssetsClient:             &AssetsClient,
		TransformsClient:         &TransformsClient,
		StreamingEndpointsClient: &StreamingEndpointsClient,
		JobsClient:               &JobsClient,
		StreamingLocatorsClient:  &StreamingLocatorsClient,
		ContentKeyPoliciesClient: &ContentKeyPoliciesClient,
		StreamingPoliciesClient:  &StreamingPoliciesClient,
		LiveEventsClient:         &LiveEventsClient,
		LiveOutputsClient:        &LiveOutputsClient,
		AssetFiltersClient:       &AssetFiltersClient,
	}
}
