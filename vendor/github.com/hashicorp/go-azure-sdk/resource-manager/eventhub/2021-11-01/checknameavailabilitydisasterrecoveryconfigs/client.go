package checknameavailabilitydisasterrecoveryconfigs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityDisasterRecoveryConfigsClient struct {
	Client *resourcemanager.Client
}

func NewCheckNameAvailabilityDisasterRecoveryConfigsClientWithBaseURI(sdkApi sdkEnv.Api) (*CheckNameAvailabilityDisasterRecoveryConfigsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "checknameavailabilitydisasterrecoveryconfigs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CheckNameAvailabilityDisasterRecoveryConfigsClient: %+v", err)
	}

	return &CheckNameAvailabilityDisasterRecoveryConfigsClient{
		Client: client,
	}, nil
}
