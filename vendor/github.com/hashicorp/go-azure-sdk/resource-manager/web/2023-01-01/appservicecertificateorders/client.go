package appservicecertificateorders

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServiceCertificateOrdersClient struct {
	Client *resourcemanager.Client
}

func NewAppServiceCertificateOrdersClientWithBaseURI(sdkApi sdkEnv.Api) (*AppServiceCertificateOrdersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "appservicecertificateorders", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AppServiceCertificateOrdersClient: %+v", err)
	}

	return &AppServiceCertificateOrdersClient{
		Client: client,
	}, nil
}
