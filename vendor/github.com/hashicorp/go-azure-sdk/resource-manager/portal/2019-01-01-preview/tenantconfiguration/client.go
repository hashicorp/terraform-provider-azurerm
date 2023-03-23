package tenantconfiguration

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantConfigurationClient struct {
	Client  autorest.Client
	baseUri string
}

func NewTenantConfigurationClientWithBaseURI(endpoint string) TenantConfigurationClient {
	return TenantConfigurationClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
