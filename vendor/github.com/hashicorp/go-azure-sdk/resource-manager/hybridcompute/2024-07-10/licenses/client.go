package licenses

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicensesClient struct {
	Client *resourcemanager.Client
}

func NewLicensesClientWithBaseURI(sdkApi sdkEnv.Api) (*LicensesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "licenses", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LicensesClient: %+v", err)
	}

	return &LicensesClient{
		Client: client,
	}, nil
}
