package managedinstanceazureadonlyauthentications

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstanceAzureADOnlyAuthenticationsClient struct {
	Client *resourcemanager.Client
}

func NewManagedInstanceAzureADOnlyAuthenticationsClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedInstanceAzureADOnlyAuthenticationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedinstanceazureadonlyauthentications", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedInstanceAzureADOnlyAuthenticationsClient: %+v", err)
	}

	return &ManagedInstanceAzureADOnlyAuthenticationsClient{
		Client: client,
	}, nil
}
