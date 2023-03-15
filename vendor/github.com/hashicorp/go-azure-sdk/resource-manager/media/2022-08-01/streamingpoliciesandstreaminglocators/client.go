package streamingpoliciesandstreaminglocators

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingPoliciesAndStreamingLocatorsClient struct {
	Client *resourcemanager.Client
}

func NewStreamingPoliciesAndStreamingLocatorsClientWithBaseURI(api environments.Api) (*StreamingPoliciesAndStreamingLocatorsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "streamingpoliciesandstreaminglocators", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StreamingPoliciesAndStreamingLocatorsClient: %+v", err)
	}

	return &StreamingPoliciesAndStreamingLocatorsClient{
		Client: client,
	}, nil
}
