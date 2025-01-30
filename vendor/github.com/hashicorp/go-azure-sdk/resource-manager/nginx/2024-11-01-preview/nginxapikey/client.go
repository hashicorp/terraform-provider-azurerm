package nginxapikey

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxApiKeyClient struct {
	Client *resourcemanager.Client
}

func NewNginxApiKeyClientWithBaseURI(sdkApi sdkEnv.Api) (*NginxApiKeyClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "nginxapikey", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NginxApiKeyClient: %+v", err)
	}

	return &NginxApiKeyClient{
		Client: client,
	}, nil
}
