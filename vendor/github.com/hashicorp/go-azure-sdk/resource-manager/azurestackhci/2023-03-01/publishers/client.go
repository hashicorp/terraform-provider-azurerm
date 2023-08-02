package publishers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublishersClient struct {
	Client *resourcemanager.Client
}

func NewPublishersClientWithBaseURI(sdkApi sdkEnv.Api) (*PublishersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "publishers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PublishersClient: %+v", err)
	}

	return &PublishersClient{
		Client: client,
	}, nil
}
