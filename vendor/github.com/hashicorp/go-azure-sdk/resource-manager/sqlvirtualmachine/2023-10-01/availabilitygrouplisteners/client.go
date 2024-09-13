package availabilitygrouplisteners

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailabilityGroupListenersClient struct {
	Client *resourcemanager.Client
}

func NewAvailabilityGroupListenersClientWithBaseURI(sdkApi sdkEnv.Api) (*AvailabilityGroupListenersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "availabilitygrouplisteners", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AvailabilityGroupListenersClient: %+v", err)
	}

	return &AvailabilityGroupListenersClient{
		Client: client,
	}, nil
}
