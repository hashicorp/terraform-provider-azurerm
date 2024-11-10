package singlesignon

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SingleSignOnClient struct {
	Client *resourcemanager.Client
}

func NewSingleSignOnClientWithBaseURI(sdkApi sdkEnv.Api) (*SingleSignOnClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "singlesignon", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SingleSignOnClient: %+v", err)
	}

	return &SingleSignOnClient{
		Client: client,
	}, nil
}
