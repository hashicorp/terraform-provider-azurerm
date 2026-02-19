package organizationresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrganizationResourcesClient struct {
	Client *resourcemanager.Client
}

func NewOrganizationResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*OrganizationResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "organizationresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OrganizationResourcesClient: %+v", err)
	}

	return &OrganizationResourcesClient{
		Client: client,
	}, nil
}
