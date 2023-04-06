package redis

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisClient struct {
	Client  autorest.Client
	baseUri string
}

func NewRedisClientWithBaseURI(endpoint string) RedisClient {
	return RedisClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
