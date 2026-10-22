package securityuserrulecollections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityUserRuleCollectionsClient struct {
	Client *resourcemanager.Client
}

func NewSecurityUserRuleCollectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*SecurityUserRuleCollectionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "securityuserrulecollections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SecurityUserRuleCollectionsClient: %+v", err)
	}

	return &SecurityUserRuleCollectionsClient{
		Client: client,
	}, nil
}
