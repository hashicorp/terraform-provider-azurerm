package contactprofile

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactProfileClient struct {
	Client  autorest.Client
	baseUri string
}

func NewContactProfileClientWithBaseURI(endpoint string) ContactProfileClient {
	return ContactProfileClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
