package redisenterprise

import "github.com/Azure/go-autorest/autorest"

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
