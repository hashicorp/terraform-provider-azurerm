package dataexport

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataExportClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDataExportClientWithBaseURI(endpoint string) DataExportClient {
	return DataExportClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
