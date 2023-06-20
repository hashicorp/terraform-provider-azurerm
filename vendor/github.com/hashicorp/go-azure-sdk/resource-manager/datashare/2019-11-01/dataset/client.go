package dataset

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataSetClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDataSetClientWithBaseURI(endpoint string) DataSetClient {
	return DataSetClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
