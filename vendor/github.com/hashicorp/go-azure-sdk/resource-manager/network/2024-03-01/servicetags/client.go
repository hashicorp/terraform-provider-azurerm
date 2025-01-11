package servicetags

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceTagsClient struct {
	Client *resourcemanager.Client
}

func NewServiceTagsClientWithBaseURI(sdkApi sdkEnv.Api) (*ServiceTagsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "servicetags", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServiceTagsClient: %+v", err)
	}

	return &ServiceTagsClient{
		Client: client,
	}, nil
}
