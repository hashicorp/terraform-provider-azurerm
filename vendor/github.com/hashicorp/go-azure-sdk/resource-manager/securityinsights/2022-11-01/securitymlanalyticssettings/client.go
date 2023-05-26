package securitymlanalyticssettings

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityMLAnalyticsSettingsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSecurityMLAnalyticsSettingsClientWithBaseURI(endpoint string) SecurityMLAnalyticsSettingsClient {
	return SecurityMLAnalyticsSettingsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
