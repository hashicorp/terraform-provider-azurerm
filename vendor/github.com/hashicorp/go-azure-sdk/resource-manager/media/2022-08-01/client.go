package v2022_08_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/accountfilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/assetsandassetfilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/contentkeypolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/liveevents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/liveoutputs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingendpoint"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingpoliciesandstreaminglocators"
)

type Client struct {
	AccountFilters                        *accountfilters.AccountFiltersClient
	AssetsAndAssetFilters                 *assetsandassetfilters.AssetsAndAssetFiltersClient
	ContentKeyPolicies                    *contentkeypolicies.ContentKeyPoliciesClient
	LiveEvents                            *liveevents.LiveEventsClient
	LiveOutputs                           *liveoutputs.LiveOutputsClient
	StreamingEndpoint                     *streamingendpoint.StreamingEndpointClient
	StreamingEndpoints                    *streamingendpoints.StreamingEndpointsClient
	StreamingPoliciesAndStreamingLocators *streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	accountFiltersClient := accountfilters.NewAccountFiltersClientWithBaseURI(endpoint)
	configureAuthFunc(&accountFiltersClient.Client)

	assetsAndAssetFiltersClient := assetsandassetfilters.NewAssetsAndAssetFiltersClientWithBaseURI(endpoint)
	configureAuthFunc(&assetsAndAssetFiltersClient.Client)

	contentKeyPoliciesClient := contentkeypolicies.NewContentKeyPoliciesClientWithBaseURI(endpoint)
	configureAuthFunc(&contentKeyPoliciesClient.Client)

	liveEventsClient := liveevents.NewLiveEventsClientWithBaseURI(endpoint)
	configureAuthFunc(&liveEventsClient.Client)

	liveOutputsClient := liveoutputs.NewLiveOutputsClientWithBaseURI(endpoint)
	configureAuthFunc(&liveOutputsClient.Client)

	streamingEndpointClient := streamingendpoint.NewStreamingEndpointClientWithBaseURI(endpoint)
	configureAuthFunc(&streamingEndpointClient.Client)

	streamingEndpointsClient := streamingendpoints.NewStreamingEndpointsClientWithBaseURI(endpoint)
	configureAuthFunc(&streamingEndpointsClient.Client)

	streamingPoliciesAndStreamingLocatorsClient := streamingpoliciesandstreaminglocators.NewStreamingPoliciesAndStreamingLocatorsClientWithBaseURI(endpoint)
	configureAuthFunc(&streamingPoliciesAndStreamingLocatorsClient.Client)

	return Client{
		AccountFilters:                        &accountFiltersClient,
		AssetsAndAssetFilters:                 &assetsAndAssetFiltersClient,
		ContentKeyPolicies:                    &contentKeyPoliciesClient,
		LiveEvents:                            &liveEventsClient,
		LiveOutputs:                           &liveOutputsClient,
		StreamingEndpoint:                     &streamingEndpointClient,
		StreamingEndpoints:                    &streamingEndpointsClient,
		StreamingPoliciesAndStreamingLocators: &streamingPoliciesAndStreamingLocatorsClient,
	}
}
