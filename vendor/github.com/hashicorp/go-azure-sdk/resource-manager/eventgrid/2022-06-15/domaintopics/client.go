package domaintopics

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DomainTopicsClient struct {
	Client *resourcemanager.Client
}

func NewDomainTopicsClientWithBaseURI(sdkApi sdkEnv.Api) (*DomainTopicsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "domaintopics", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DomainTopicsClient: %+v", err)
	}

	return &DomainTopicsClient{
		Client: client,
	}, nil
}
