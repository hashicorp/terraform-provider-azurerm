package dnsprivateviews

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsPrivateViewsClient struct {
	Client *resourcemanager.Client
}

func NewDnsPrivateViewsClientWithBaseURI(sdkApi sdkEnv.Api) (*DnsPrivateViewsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dnsprivateviews", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DnsPrivateViewsClient: %+v", err)
	}

	return &DnsPrivateViewsClient{
		Client: client,
	}, nil
}
