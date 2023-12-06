package adminrulecollections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdminRuleCollectionsClient struct {
	Client *resourcemanager.Client
}

func NewAdminRuleCollectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*AdminRuleCollectionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "adminrulecollections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AdminRuleCollectionsClient: %+v", err)
	}

	return &AdminRuleCollectionsClient{
		Client: client,
	}, nil
}
