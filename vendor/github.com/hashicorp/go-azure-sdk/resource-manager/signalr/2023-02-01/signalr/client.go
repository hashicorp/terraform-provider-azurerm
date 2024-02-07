package signalr

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SignalRClient struct {
	Client *resourcemanager.Client
}

func NewSignalRClientWithBaseURI(sdkApi sdkEnv.Api) (*SignalRClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "signalr", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SignalRClient: %+v", err)
	}

	return &SignalRClient{
		Client: client,
	}, nil
}
