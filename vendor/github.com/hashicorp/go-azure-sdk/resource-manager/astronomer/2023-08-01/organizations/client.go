package organizations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrganizationsClient struct {
	Client *resourcemanager.Client
}

func NewOrganizationsClientWithBaseURI(sdkApi sdkEnv.Api) (*OrganizationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "organizations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OrganizationsClient: %+v", err)
	}

	return &OrganizationsClient{
		Client: client,
	}, nil
}
