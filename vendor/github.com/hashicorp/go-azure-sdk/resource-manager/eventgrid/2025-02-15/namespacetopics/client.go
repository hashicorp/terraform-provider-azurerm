package namespacetopics

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespaceTopicsClient struct {
	Client *resourcemanager.Client
}

func NewNamespaceTopicsClientWithBaseURI(sdkApi sdkEnv.Api) (*NamespaceTopicsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "namespacetopics", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NamespaceTopicsClient: %+v", err)
	}

	return &NamespaceTopicsClient{
		Client: client,
	}, nil
}
