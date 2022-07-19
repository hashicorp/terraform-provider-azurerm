package checknameavailabilitydisasterrecoveryconfigs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityDisasterRecoveryConfigsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCheckNameAvailabilityDisasterRecoveryConfigsClientWithBaseURI(endpoint string) CheckNameAvailabilityDisasterRecoveryConfigsClient {
	return CheckNameAvailabilityDisasterRecoveryConfigsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
