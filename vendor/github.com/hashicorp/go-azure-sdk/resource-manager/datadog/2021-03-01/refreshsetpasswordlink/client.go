package refreshsetpasswordlink

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RefreshSetPasswordLinkClient struct {
	Client *resourcemanager.Client
}

func NewRefreshSetPasswordLinkClientWithBaseURI(sdkApi sdkEnv.Api) (*RefreshSetPasswordLinkClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "refreshsetpasswordlink", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RefreshSetPasswordLinkClient: %+v", err)
	}

	return &RefreshSetPasswordLinkClient{
		Client: client,
	}, nil
}
