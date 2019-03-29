package capi

import (
	"github.com/Azure/go-autorest/autorest"
	"strings"
)

type BaseClient struct {
	autorest.Client
	ID string //client identifier for errors

	AccountName string
	AccountKey  AccountKey
	BaseURI     string

	Version string //optional
}

func NewClient(clientId, cosmosAccountName string, cosmosAccountKey AccountKey, cosmosURLSuffix, version string) BaseClient {
	return BaseClient{
		Client:      autorest.NewClientWithUserAgent("go-cosmos-sdk.0.0.0.0"),
		ID:          clientId,
		AccountName: cosmosAccountName,
		AccountKey:  cosmosAccountKey,
		BaseURI:     "https://" + cosmosAccountName + "." + cosmosURLSuffix,
		Version:     version,
	}
}

type MakesResourcePath interface {
	GenerateResourcePath(map[string]interface{}) string
}

func (c BaseClient) GenerateResourcePath(path string, parameters map[string]string) string {
	for key, value := range parameters {
		path = strings.Replace(path, "{"+key+"}", value, -1)
	}
	return path
}
