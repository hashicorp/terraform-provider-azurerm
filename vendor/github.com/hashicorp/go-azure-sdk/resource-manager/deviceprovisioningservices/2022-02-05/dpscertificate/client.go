package dpscertificate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DpsCertificateClient struct {
	Client *resourcemanager.Client
}

func NewDpsCertificateClientWithBaseURI(sdkApi sdkEnv.Api) (*DpsCertificateClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dpscertificate", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DpsCertificateClient: %+v", err)
	}

	return &DpsCertificateClient{
		Client: client,
	}, nil
}
