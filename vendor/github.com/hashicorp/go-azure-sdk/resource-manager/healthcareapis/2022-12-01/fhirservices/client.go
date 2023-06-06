package fhirservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FhirServicesClient struct {
	Client *resourcemanager.Client
}

func NewFhirServicesClientWithBaseURI(api environments.Api) (*FhirServicesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "fhirservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FhirServicesClient: %+v", err)
	}

	return &FhirServicesClient{
		Client: client,
	}, nil
}
