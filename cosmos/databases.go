package cosmos

import (
	"context"
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/cosmos/capi"
	"strconv"
)

// Client
type DatabasesClient struct {
	capi.BaseClient
}

func NewDatabasesClientWithDefaults(cosmosAccountName string, cosmosAccountKey capi.AccountKey) DatabasesClient {
	return NewDatabasesClient(cosmosAccountName, cosmosAccountKey, capi.DefaultURLSuffix, capi.DefaultCosmosAPIVersion)
}

func NewDatabasesClient(cosmosAccountName string, cosmosAccountKey capi.AccountKey, cosmosURLSuffix, version string) DatabasesClient {
	return DatabasesClient{
		capi.NewClient("cosmos.DatabasesClient", cosmosAccountName, cosmosAccountKey, cosmosURLSuffix, version),
	}
}

// Object
type Database struct {
	capi.BaseResource

	OfferThroughput *int `json:"-"`
}

// PathBase
type PathToDatabase struct {
	capi.PathBase

	Database string
}

func ParseDatabasePath(path string) (PathToDatabase, error) {
	p, err := capi.ParsePath(path)
	if err != nil {
		return PathToDatabase{}, fmt.Errorf("Unable to parse database path: %v", err)
	}

	dbp := PathToDatabase{PathBase: p}

	if v, ok := p.Parts["dbs"]; ok {
		dbp.Database = v
		return dbp, nil
	}

	return dbp, fmt.Errorf("Unable to parse database path, missing `dbs` component")
}

//should this retirn a PathBase object?
func BuildDatabasePath(databaseName string) PathToDatabase {
	p, _ := ParseDatabasePath("/dbs/" + databaseName) //cannot.. well should not fail
	return p
}

// Methods
func (c DatabasesClient) Create(ctx context.Context, databaseName string, input Database) (result Database, err error) {
	path := BuildDatabasePath(databaseName)
	input.ID = &databaseName

	preparers := []autorest.PrepareDecorator{
		autorest.WithJSON(input),
	}

	if input.OfferThroughput != nil {
		preparers = append(preparers, autorest.WithHeader("x-ms-offer-throughput", strconv.Itoa(*input.OfferThroughput)))
	}

	resp, err := c.BaseClient.Create(ctx, path.GetCreatePath(), preparers, []autorest.RespondDecorator{autorest.ByUnmarshallingJSON(&result)})
	result.PopulateBase(err, resp)
	return result, err
}

func (c DatabasesClient) Get(ctx context.Context, databaseName string) (result Database, err error) {
	path := BuildDatabasePath(databaseName)

	resp, err := c.BaseClient.Get(ctx, path.String, autorest.ByUnmarshallingJSON(&result))
	result.PopulateBase(err, resp)
	return result, err
}

func (c DatabasesClient) Delete(ctx context.Context, databaseName string) (result *autorest.Response, err error) {
	return c.BaseClient.Delete(ctx, BuildDatabasePath(databaseName).String)
}
