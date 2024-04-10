package diagnosticsettingscategories

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticSettingsCategoriesClient struct {
	Client *resourcemanager.Client
}

func NewDiagnosticSettingsCategoriesClientWithBaseURI(sdkApi sdkEnv.Api) (*DiagnosticSettingsCategoriesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "diagnosticsettingscategories", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DiagnosticSettingsCategoriesClient: %+v", err)
	}

	return &DiagnosticSettingsCategoriesClient{
		Client: client,
	}, nil
}
