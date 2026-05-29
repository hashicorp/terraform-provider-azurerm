package assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssignmentClient struct {
	Client *resourcemanager.Client
}

func NewAssignmentClientWithBaseURI(sdkApi sdkEnv.Api) (*AssignmentClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "assignment", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AssignmentClient: %+v", err)
	}

	return &AssignmentClient{
		Client: client,
	}, nil
}
