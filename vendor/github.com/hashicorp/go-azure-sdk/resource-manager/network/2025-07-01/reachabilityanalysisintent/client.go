package reachabilityanalysisintent

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReachabilityAnalysisIntentClient struct {
	Client *resourcemanager.Client
}

func NewReachabilityAnalysisIntentClientWithBaseURI(sdkApi sdkEnv.Api) (*ReachabilityAnalysisIntentClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "reachabilityanalysisintent", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ReachabilityAnalysisIntentClient: %+v", err)
	}

	return &ReachabilityAnalysisIntentClient{
		Client: client,
	}, nil
}
