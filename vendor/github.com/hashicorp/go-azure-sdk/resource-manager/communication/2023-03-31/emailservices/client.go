package emailservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EmailServicesClient struct {
	Client *resourcemanager.Client
}

func NewEmailServicesClientWithBaseURI(api environments.Api) (*EmailServicesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "emailservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EmailServicesClient: %+v", err)
	}

	return &EmailServicesClient{
		Client: client,
	}, nil
}
