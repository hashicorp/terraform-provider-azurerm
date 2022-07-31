package recordsets

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecordSetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRecordSetsClientWithBaseURI(endpoint string) RecordSetsClient {
	return RecordSetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
