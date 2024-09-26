package azuresdkhacks

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest"
)

// TODO 4.0: check if it could be removed on 4.0
// workaround for https://github.com/Azure/azure-rest-api-specs/issues/22572
// the swagger lack definition of `certificateCreateOptions`.

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

const defaultApiVersion = "2024-04-01"

func userAgent() string {
	return fmt.Sprintf("hashicorp/go-azure-sdk/vaultcertificates/%s", defaultApiVersion)
}

type VaultCertificatesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVaultCertificatesClientWithBaseURI(endpoint string) VaultCertificatesClient {
	return VaultCertificatesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
