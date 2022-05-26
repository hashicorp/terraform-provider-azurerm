package fluxconfigurationoperationstatus

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluxConfigurationOperationStatusClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFluxConfigurationOperationStatusClientWithBaseURI(endpoint string) FluxConfigurationOperationStatusClient {
	return FluxConfigurationOperationStatusClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
