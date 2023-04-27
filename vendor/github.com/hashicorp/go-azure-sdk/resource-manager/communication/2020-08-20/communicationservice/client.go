package communicationservice

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommunicationServiceClient struct {
	Client *resourcemanager.Client
}

func NewCommunicationServiceClientWithBaseURI(api environments.Api) (*CommunicationServiceClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "communicationservice", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CommunicationServiceClient: %+v", err)
	}

	return &CommunicationServiceClient{
		Client: client,
	}, nil
}
