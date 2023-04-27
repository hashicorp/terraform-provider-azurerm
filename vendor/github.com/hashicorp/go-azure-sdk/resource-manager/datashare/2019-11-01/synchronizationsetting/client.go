package synchronizationsetting

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynchronizationSettingClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSynchronizationSettingClientWithBaseURI(endpoint string) SynchronizationSettingClient {
	return SynchronizationSettingClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
