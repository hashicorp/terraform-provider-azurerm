package registrationdefinitions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistrationDefinitionsClient struct {
	Client *resourcemanager.Client
}

func NewRegistrationDefinitionsClientWithBaseURI(sdkApi sdkEnv.Api) (*RegistrationDefinitionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "registrationdefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RegistrationDefinitionsClient: %+v", err)
	}

	return &RegistrationDefinitionsClient{
		Client: client,
	}, nil
}
