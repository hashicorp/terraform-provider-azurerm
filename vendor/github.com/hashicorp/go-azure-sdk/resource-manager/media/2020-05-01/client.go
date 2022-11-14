package v2020_05_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/accountfilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/accounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/assetsandassetfilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/contentkeypolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/encodings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/liveevents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/liveoutputs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/streamingendpoint"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/streamingendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/streamingpoliciesandstreaminglocators"
)

type Client struct {
	AccountFilters                        *accountfilters.AccountFiltersClient
	Accounts                              *accounts.AccountsClient
	AssetsAndAssetFilters                 *assetsandassetfilters.AssetsAndAssetFiltersClient
	ContentKeyPolicies                    *contentkeypolicies.ContentKeyPoliciesClient
	Encodings                             *encodings.EncodingsClient
	LiveEvents                            *liveevents.LiveEventsClient
	LiveOutputs                           *liveoutputs.LiveOutputsClient
	StreamingEndpoint                     *streamingendpoint.StreamingEndpointClient
	StreamingEndpoints                    *streamingendpoints.StreamingEndpointsClient
	StreamingPoliciesAndStreamingLocators *streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	accountFiltersClient := accountfilters.NewAccountFiltersClientWithBaseURI(endpoint)
	configureAuthFunc(&accountFiltersClient.Client)

	accountsClient := accounts.NewAccountsClientWithBaseURI(endpoint)
	configureAuthFunc(&accountsClient.Client)

	assetsAndAssetFiltersClient := assetsandassetfilters.NewAssetsAndAssetFiltersClientWithBaseURI(endpoint)
	configureAuthFunc(&assetsAndAssetFiltersClient.Client)

	contentKeyPoliciesClient := contentkeypolicies.NewContentKeyPoliciesClientWithBaseURI(endpoint)
	configureAuthFunc(&contentKeyPoliciesClient.Client)

	encodingsClient := encodings.NewEncodingsClientWithBaseURI(endpoint)
	configureAuthFunc(&encodingsClient.Client)

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
		Accounts:                              &accountsClient,
		AssetsAndAssetFilters:                 &assetsAndAssetFiltersClient,
		ContentKeyPolicies:                    &contentKeyPoliciesClient,
		Encodings:                             &encodingsClient,
		LiveEvents:                            &liveEventsClient,
		LiveOutputs:                           &liveOutputsClient,
		StreamingEndpoint:                     &streamingEndpointClient,
		StreamingEndpoints:                    &streamingEndpointsClient,
		StreamingPoliciesAndStreamingLocators: &streamingPoliciesAndStreamingLocatorsClient,
	}
}
