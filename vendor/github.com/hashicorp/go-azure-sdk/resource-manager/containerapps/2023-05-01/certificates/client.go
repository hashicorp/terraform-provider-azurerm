package certificates

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatesClient struct {
	Client *resourcemanager.Client
}

func NewCertificatesClientWithBaseURI(sdkApi sdkEnv.Api) (*CertificatesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "certificates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CertificatesClient: %+v", err)
	}

	return &CertificatesClient{
		Client: client,
	}, nil
}
