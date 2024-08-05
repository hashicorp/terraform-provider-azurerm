package certificateobjectlocalrulestack

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateObjectLocalRulestackClient struct {
	Client *resourcemanager.Client
}

func NewCertificateObjectLocalRulestackClientWithBaseURI(sdkApi sdkEnv.Api) (*CertificateObjectLocalRulestackClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "certificateobjectlocalrulestack", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CertificateObjectLocalRulestackClient: %+v", err)
	}

	return &CertificateObjectLocalRulestackClient{
		Client: client,
	}, nil
}
