package certificateobjectglobalrulestack

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateObjectGlobalRulestackClient struct {
	Client *resourcemanager.Client
}

func NewCertificateObjectGlobalRulestackClientWithBaseURI(sdkApi sdkEnv.Api) (*CertificateObjectGlobalRulestackClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "certificateobjectglobalrulestack", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CertificateObjectGlobalRulestackClient: %+v", err)
	}

	return &CertificateObjectGlobalRulestackClient{
		Client: client,
	}, nil
}
