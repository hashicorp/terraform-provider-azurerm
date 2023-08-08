package v2021_11_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/accountfilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/accounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/assetsandassetfilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/contentkeypolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/encodings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/liveevents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/liveoutputs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/streamingendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/streamingpoliciesandstreaminglocators"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AccountFilters                        *accountfilters.AccountFiltersClient
	Accounts                              *accounts.AccountsClient
	AssetsAndAssetFilters                 *assetsandassetfilters.AssetsAndAssetFiltersClient
	ContentKeyPolicies                    *contentkeypolicies.ContentKeyPoliciesClient
	Encodings                             *encodings.EncodingsClient
	LiveEvents                            *liveevents.LiveEventsClient
	LiveOutputs                           *liveoutputs.LiveOutputsClient
	StreamingEndpoints                    *streamingendpoints.StreamingEndpointsClient
	StreamingPoliciesAndStreamingLocators *streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	accountFiltersClient, err := accountfilters.NewAccountFiltersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AccountFilters client: %+v", err)
	}
	configureFunc(accountFiltersClient.Client)

	accountsClient, err := accounts.NewAccountsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Accounts client: %+v", err)
	}
	configureFunc(accountsClient.Client)

	assetsAndAssetFiltersClient, err := assetsandassetfilters.NewAssetsAndAssetFiltersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AssetsAndAssetFilters client: %+v", err)
	}
	configureFunc(assetsAndAssetFiltersClient.Client)

	contentKeyPoliciesClient, err := contentkeypolicies.NewContentKeyPoliciesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ContentKeyPolicies client: %+v", err)
	}
	configureFunc(contentKeyPoliciesClient.Client)

	encodingsClient, err := encodings.NewEncodingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Encodings client: %+v", err)
	}
	configureFunc(encodingsClient.Client)

	liveEventsClient, err := liveevents.NewLiveEventsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LiveEvents client: %+v", err)
	}
	configureFunc(liveEventsClient.Client)

	liveOutputsClient, err := liveoutputs.NewLiveOutputsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LiveOutputs client: %+v", err)
	}
	configureFunc(liveOutputsClient.Client)

	streamingEndpointsClient, err := streamingendpoints.NewStreamingEndpointsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StreamingEndpoints client: %+v", err)
	}
	configureFunc(streamingEndpointsClient.Client)

	streamingPoliciesAndStreamingLocatorsClient, err := streamingpoliciesandstreaminglocators.NewStreamingPoliciesAndStreamingLocatorsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StreamingPoliciesAndStreamingLocators client: %+v", err)
	}
	configureFunc(streamingPoliciesAndStreamingLocatorsClient.Client)

	return &Client{
		AccountFilters:                        accountFiltersClient,
		Accounts:                              accountsClient,
		AssetsAndAssetFilters:                 assetsAndAssetFiltersClient,
		ContentKeyPolicies:                    contentKeyPoliciesClient,
		Encodings:                             encodingsClient,
		LiveEvents:                            liveEventsClient,
		LiveOutputs:                           liveOutputsClient,
		StreamingEndpoints:                    streamingEndpointsClient,
		StreamingPoliciesAndStreamingLocators: streamingPoliciesAndStreamingLocatorsClient,
	}, nil
}
