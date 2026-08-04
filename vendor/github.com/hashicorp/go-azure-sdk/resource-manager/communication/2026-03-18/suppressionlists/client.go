package suppressionlists

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SuppressionListsClient struct {
	Client *resourcemanager.Client
}

func NewSuppressionListsClientWithBaseURI(sdkApi sdkEnv.Api) (*SuppressionListsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "suppressionlists", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SuppressionListsClient: %+v", err)
	}

	return &SuppressionListsClient{
		Client: client,
	}, nil
}
