package gatewaylisttrace

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayListTraceClient struct {
	Client *resourcemanager.Client
}

func NewGatewayListTraceClientWithBaseURI(sdkApi sdkEnv.Api) (*GatewayListTraceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "gatewaylisttrace", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GatewayListTraceClient: %+v", err)
	}

	return &GatewayListTraceClient{
		Client: client,
	}, nil
}
