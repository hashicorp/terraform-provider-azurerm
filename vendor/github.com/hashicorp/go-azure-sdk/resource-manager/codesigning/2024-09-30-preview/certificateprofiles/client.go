package certificateprofiles

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateProfilesClient struct {
	Client *resourcemanager.Client
}

func NewCertificateProfilesClientWithBaseURI(sdkApi sdkEnv.Api) (*CertificateProfilesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "certificateprofiles", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CertificateProfilesClient: %+v", err)
	}

	return &CertificateProfilesClient{
		Client: client,
	}, nil
}
