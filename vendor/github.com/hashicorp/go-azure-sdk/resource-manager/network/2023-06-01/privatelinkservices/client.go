package privatelinkservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServicesClient struct {
	Client *resourcemanager.Client
}

func NewPrivateLinkServicesClientWithBaseURI(sdkApi sdkEnv.Api) (*PrivateLinkServicesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "privatelinkservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrivateLinkServicesClient: %+v", err)
	}

	return &PrivateLinkServicesClient{
		Client: client,
	}, nil
}
