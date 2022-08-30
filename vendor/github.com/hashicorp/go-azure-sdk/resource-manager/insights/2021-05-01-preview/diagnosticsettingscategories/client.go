package diagnosticsettingscategories

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticSettingsCategoriesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDiagnosticSettingsCategoriesClientWithBaseURI(endpoint string) DiagnosticSettingsCategoriesClient {
	return DiagnosticSettingsCategoriesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
