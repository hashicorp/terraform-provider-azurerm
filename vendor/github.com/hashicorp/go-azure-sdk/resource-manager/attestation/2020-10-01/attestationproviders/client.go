package attestationproviders

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttestationProvidersClient struct {
	Client *resourcemanager.Client
}

func NewAttestationProvidersClientWithBaseURI(sdkApi sdkEnv.Api) (*AttestationProvidersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "attestationproviders", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AttestationProvidersClient: %+v", err)
	}

	return &AttestationProvidersClient{
		Client: client,
	}, nil
}
