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

/*
func ByUnmarshallingDatabase(db *Database) autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {

			err := r.Respond(resp)

			if err == nil {
				b, errInner := ioutil.ReadAll(resp.Body)
				// Some responses might include a BOM, remove for successful unmarshalling
				b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
				if errInner != nil {
					err = fmt.Errorf("Error occurred reading http.Response#Body - Error = '%v'", errInner)
				} else if len(strings.Trim(string(b), " ")) > 0 {
					errInner = json.Unmarshal(b, v)
					if errInner != nil {
						err = fmt.Errorf("Error occurred unmarshalling JSON - Error = '%v' JSON = '%s'", errInner, string(b))
					}
				}
			}


			return err
		})
	}
}*/

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
func GenerateDatabasePath(databaseName string) PathToDatabase {
	p, _ := ParseDatabasePath("/dbs/" + databaseName) //cannot.. well should not fail
	return p
}

// Methods
func (c DatabasesClient) Create(ctx context.Context, databaseName string, db Database) (result Database, err error) {
	path := GenerateDatabasePath(databaseName)
	db.ID = &databaseName

	preparers := []autorest.PrepareDecorator{
		autorest.WithJSON(db),
	}

	if db.OfferThroughput != nil {
		preparers = append(preparers, autorest.WithHeader("x-ms-offer-throughput", strconv.Itoa(*db.OfferThroughput)))
	}

	resp, err := c.BaseClient.Create(ctx, path.GetCreatePath(), preparers, []autorest.RespondDecorator{autorest.ByUnmarshallingJSON(&result)})

	result.Response = *resp
	if err == nil {
		v := resp.Header.Get("x-ms-offer-throughput")

		if v != "" {
			//populate
			i, err := strconv.Atoi(resp.Header.Get("x-ms-offer-throughput"))
			if err != nil {
				return result, fmt.Errorf("unable to")
			}
			result.OfferThroughput = &i
		}
	}

	return result, err
}

func (c DatabasesClient) Get(ctx context.Context, databaseName string) (result Database, err error) {
	path := GenerateDatabasePath(databaseName)

	resp, err := c.BaseClient.Get(ctx, path.Path, autorest.ByUnmarshallingJSON(&result))

	result.Response = *resp
	if err == nil {
		v := resp.Header.Get("x-ms-offer-throughput")

		if v != "" {
			//populate
			i, err := strconv.Atoi(resp.Header.Get("x-ms-offer-throughput"))
			if err != nil {
				return result, fmt.Errorf("unable to")
			}
			result.OfferThroughput = &i
		}
	}

	return result, err
}

func (c DatabasesClient) Delete(ctx context.Context, databaseName string) (result *autorest.Response, err error) {
	return c.BaseClient.Delete(ctx, GenerateDatabasePath(databaseName).Path)
}
