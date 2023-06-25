package softwareupdateconfiguration

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSoftwareUpdateConfigurationClientWithBaseURI(endpoint string) SoftwareUpdateConfigurationClient {
	return SoftwareUpdateConfigurationClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
