package systemtopics

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemTopicsClient struct {
	Client *resourcemanager.Client
}

func NewSystemTopicsClientWithBaseURI(sdkApi sdkEnv.Api) (*SystemTopicsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "systemtopics", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SystemTopicsClient: %+v", err)
	}

	return &SystemTopicsClient{
		Client: client,
	}, nil
}
