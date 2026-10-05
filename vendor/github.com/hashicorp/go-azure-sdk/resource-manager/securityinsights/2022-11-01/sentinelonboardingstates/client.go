package sentinelonboardingstates

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SentinelOnboardingStatesClient struct {
	Client *resourcemanager.Client
}

func NewSentinelOnboardingStatesClientWithBaseURI(sdkApi sdkEnv.Api) (*SentinelOnboardingStatesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "sentinelonboardingstates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SentinelOnboardingStatesClient: %+v", err)
	}

	return &SentinelOnboardingStatesClient{
		Client: client,
	}, nil
}
