package checknameavailability

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityClient struct {
	Client *resourcemanager.Client
}

func NewCheckNameAvailabilityClientWithBaseURI(sdkApi sdkEnv.Api) (*CheckNameAvailabilityClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "checknameavailability", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CheckNameAvailabilityClient: %+v", err)
	}

	return &CheckNameAvailabilityClient{
		Client: client,
	}, nil
}
