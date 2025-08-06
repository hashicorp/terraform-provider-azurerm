package componentapikeysapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComponentApiKeysAPIsClient struct {
	Client *resourcemanager.Client
}

func NewComponentApiKeysAPIsClientWithBaseURI(sdkApi sdkEnv.Api) (*ComponentApiKeysAPIsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "componentapikeysapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ComponentApiKeysAPIsClient: %+v", err)
	}

	return &ComponentApiKeysAPIsClient{
		Client: client,
	}, nil
}
