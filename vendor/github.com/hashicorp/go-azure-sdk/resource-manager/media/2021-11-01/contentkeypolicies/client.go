package contentkeypolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewContentKeyPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*ContentKeyPoliciesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "contentkeypolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ContentKeyPoliciesClient: %+v", err)
	}

	return &ContentKeyPoliciesClient{
		Client: client,
	}, nil
}
