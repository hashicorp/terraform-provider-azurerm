package nginxcertificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxCertificateClient struct {
	Client *resourcemanager.Client
}

func NewNginxCertificateClientWithBaseURI(sdkApi sdkEnv.Api) (*NginxCertificateClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "nginxcertificate", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NginxCertificateClient: %+v", err)
	}

	return &NginxCertificateClient{
		Client: client,
	}, nil
}
