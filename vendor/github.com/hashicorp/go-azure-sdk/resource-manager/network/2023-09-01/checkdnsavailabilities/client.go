package checkdnsavailabilities

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckDnsAvailabilitiesClient struct {
	Client *resourcemanager.Client
}

func NewCheckDnsAvailabilitiesClientWithBaseURI(sdkApi sdkEnv.Api) (*CheckDnsAvailabilitiesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "checkdnsavailabilities", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CheckDnsAvailabilitiesClient: %+v", err)
	}

	return &CheckDnsAvailabilitiesClient{
		Client: client,
	}, nil
}
