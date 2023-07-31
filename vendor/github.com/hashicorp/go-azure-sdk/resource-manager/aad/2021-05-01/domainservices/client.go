package domainservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainServicesClient struct {
	Client *resourcemanager.Client
}

func NewDomainServicesClientWithBaseURI(api environments.Api) (*DomainServicesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "domainservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DomainServicesClient: %+v", err)
	}

	return &DomainServicesClient{
		Client: client,
	}, nil
}
