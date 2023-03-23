package disasterrecoveryconfigs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DisasterRecoveryConfigsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDisasterRecoveryConfigsClientWithBaseURI(endpoint string) DisasterRecoveryConfigsClient {
	return DisasterRecoveryConfigsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
