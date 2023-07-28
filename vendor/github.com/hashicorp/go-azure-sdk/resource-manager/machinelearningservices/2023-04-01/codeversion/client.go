package codeversion

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CodeVersionClient struct {
	Client *resourcemanager.Client
}

func NewCodeVersionClientWithBaseURI(api environments.Api) (*CodeVersionClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "codeversion", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CodeVersionClient: %+v", err)
	}

	return &CodeVersionClient{
		Client: client,
	}, nil
}
