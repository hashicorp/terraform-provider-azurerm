package enrollmentaccounts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnrollmentAccountsClient struct {
	Client *resourcemanager.Client
}

func NewEnrollmentAccountsClientWithBaseURI(sdkApi sdkEnv.Api) (*EnrollmentAccountsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "enrollmentaccounts", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EnrollmentAccountsClient: %+v", err)
	}

	return &EnrollmentAccountsClient{
		Client: client,
	}, nil
}
