package redisenterprise

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisEnterpriseClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRedisEnterpriseClientWithBaseURI(endpoint string) RedisEnterpriseClient {
	return RedisEnterpriseClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
