package vaultcertificates

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VaultCertificatesClient struct {
	Client *resourcemanager.Client
}

func NewVaultCertificatesClientWithBaseURI(sdkApi sdkEnv.Api) (*VaultCertificatesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "vaultcertificates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VaultCertificatesClient: %+v", err)
	}

	return &VaultCertificatesClient{
		Client: client,
	}, nil
}
