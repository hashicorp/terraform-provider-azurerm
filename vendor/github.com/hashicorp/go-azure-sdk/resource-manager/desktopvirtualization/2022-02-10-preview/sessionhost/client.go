package sessionhost

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SessionHostClient struct {
	Client *resourcemanager.Client
}

func NewSessionHostClientWithBaseURI(sdkApi sdkEnv.Api) (*SessionHostClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "sessionhost", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SessionHostClient: %+v", err)
	}

	return &SessionHostClient{
		Client: client,
	}, nil
}
