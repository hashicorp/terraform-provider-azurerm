package diskencryptionsets

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskEncryptionSetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDiskEncryptionSetsClientWithBaseURI(endpoint string) DiskEncryptionSetsClient {
	return DiskEncryptionSetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
