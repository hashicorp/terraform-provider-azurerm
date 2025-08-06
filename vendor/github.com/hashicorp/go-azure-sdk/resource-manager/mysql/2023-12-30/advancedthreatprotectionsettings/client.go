package advancedthreatprotectionsettings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvancedThreatProtectionSettingsClient struct {
	Client *resourcemanager.Client
}

func NewAdvancedThreatProtectionSettingsClientWithBaseURI(sdkApi sdkEnv.Api) (*AdvancedThreatProtectionSettingsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "advancedthreatprotectionsettings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AdvancedThreatProtectionSettingsClient: %+v", err)
	}

	return &AdvancedThreatProtectionSettingsClient{
		Client: client,
	}, nil
}
