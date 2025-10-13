package managedenvironments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedEnvironmentsClient struct {
	Client *resourcemanager.Client
}

func NewManagedEnvironmentsClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedEnvironmentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedenvironments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedEnvironmentsClient: %+v", err)
	}

	return &ManagedEnvironmentsClient{
		Client: client,
	}, nil
}
