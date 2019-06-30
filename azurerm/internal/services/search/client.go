package search

import "github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"

type Client struct {
	AdminKeysClient search.AdminKeysClient
	ServicesClient  search.ServicesClient
}
