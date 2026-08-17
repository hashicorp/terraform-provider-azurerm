package v2025_10_13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	CertificateProfiles *certificateprofiles.CertificateProfilesClient
	CodeSigningAccounts *codesigningaccounts.CodeSigningAccountsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	certificateProfilesClient, err := certificateprofiles.NewCertificateProfilesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CertificateProfiles client: %+v", err)
	}
	configureFunc(certificateProfilesClient.Client)

	codeSigningAccountsClient, err := codesigningaccounts.NewCodeSigningAccountsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CodeSigningAccounts client: %+v", err)
	}
	configureFunc(codeSigningAccountsClient.Client)

	return &Client{
		CertificateProfiles: certificateProfilesClient,
		CodeSigningAccounts: codeSigningAccountsClient,
	}, nil
}
