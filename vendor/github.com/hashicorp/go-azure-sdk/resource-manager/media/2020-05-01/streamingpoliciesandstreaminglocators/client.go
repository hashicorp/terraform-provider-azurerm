package streamingpoliciesandstreaminglocators

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingPoliciesAndStreamingLocatorsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewStreamingPoliciesAndStreamingLocatorsClientWithBaseURI(endpoint string) StreamingPoliciesAndStreamingLocatorsClient {
	return StreamingPoliciesAndStreamingLocatorsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
