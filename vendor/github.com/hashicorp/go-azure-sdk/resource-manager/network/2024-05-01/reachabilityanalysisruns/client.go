package reachabilityanalysisruns

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReachabilityAnalysisRunsClient struct {
	Client *resourcemanager.Client
}

func NewReachabilityAnalysisRunsClientWithBaseURI(sdkApi sdkEnv.Api) (*ReachabilityAnalysisRunsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "reachabilityanalysisruns", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ReachabilityAnalysisRunsClient: %+v", err)
	}

	return &ReachabilityAnalysisRunsClient{
		Client: client,
	}, nil
}
