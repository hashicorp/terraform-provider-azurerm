package dppfeaturesupport

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DppFeatureSupportClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDppFeatureSupportClientWithBaseURI(endpoint string) DppFeatureSupportClient {
	return DppFeatureSupportClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
