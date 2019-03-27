package cosmos

import (
	"github.com/Azure/go-autorest/autorest"
)

const DefaultCosmosURLSuffix = "documents.azure.com"

type AccountKey struct {
	Type    string
	Version string
	Key     string
}

func NewAccountKeyWithDefaults(key string) AccountKey {
	return NewAccountKey(key, "master", "1.0")
}

func NewAccountKey(key, keyType, keyVersion string) AccountKey {
	return AccountKey{
		Type:    keyType,
		Version: keyVersion,
		Key:     key,
	}
}

type BaseClient struct {
	autorest.Client
	ID string //client identifier for errors

	AccountName string
	AccountKey  AccountKey
	BaseURI     string

	Version string //optional
}

func newClient(clientId, cosmosAccountName string, cosmosAccountKey AccountKey, cosmosURLSuffix, version string) BaseClient {
	return BaseClient{
		Client:      autorest.NewClientWithUserAgent("go-cosmos-sdk.0.0.0.0"),
		ID:          clientId,
		AccountName: cosmosAccountName,
		AccountKey:  cosmosAccountKey,
		BaseURI:     "https://" + cosmosAccountName + "." + cosmosURLSuffix,
		Version:     version,
	}
}
