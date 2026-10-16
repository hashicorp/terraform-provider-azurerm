package nginxcertificates

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxCertificatesClient struct {
	Client *resourcemanager.Client
}

func NewNginxCertificatesClientWithBaseURI(sdkApi sdkEnv.Api) (*NginxCertificatesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "nginxcertificates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NginxCertificatesClient: %+v", err)
	}

	return &NginxCertificatesClient{
		Client: client,
	}, nil
}
