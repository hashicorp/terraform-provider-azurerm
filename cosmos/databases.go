package cosmos

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"net/http"
)

// Client Definitions
type DatabasesClient struct {
	BaseClient
}

func NewDatabasesClient(cosmosAccountName string, cosmosAccountKey AccountKey) DatabasesClient {
	return NewDatabasesClientWithSuffix(cosmosAccountName, cosmosAccountKey, DefaultCosmosURLSuffix)
}

func NewDatabasesClientWithSuffix(cosmosAccountName string, cosmosAccountKey AccountKey, cosmosURLSuffix string) DatabasesClient {
	return NewDatabasesClientWithVersion(cosmosAccountName, cosmosAccountKey, cosmosURLSuffix, "2017-02-22")
}

func NewDatabasesClientWithVersion(cosmosAccountName string, cosmosAccountKey AccountKey, cosmosURLSuffix, version string) DatabasesClient {
	return DatabasesClient{
		newClient("cosmos.DatabasesClient", cosmosAccountName, cosmosAccountKey, cosmosURLSuffix, version),
	}
}

// Objects
type Database struct {
	APIResponse

	//throughput
}

func ByPopulatingEntity(e APIResponse) autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			err := r.Respond(resp)
			if err != nil {
				return err
			}

			e.Response = autorest.Response{resp}

			return nil
		})
	}
}

// Methods
func (client DatabasesClient) Create(ctx context.Context, databaseName string) (result Database, err error) {

	db := Database{
		APIResponse: APIResponse{
			ID: &databaseName,
		},
	}

	req, err := client.Preparer(ctx, "Create", http.MethodPost, "/dbs", map[string]interface{}{}, autorest.WithJSON(db))
	if err != nil {
		return result, err
	}

	resp, err := client.Sender("Create", req)
	if err == nil {
		resp, err = client.Responder("Create", resp.Response, http.StatusCreated, autorest.ByUnmarshallingJSON(&result))
		if err == nil {
			result.Path = "cosmos/dbs/" + databaseName
		}
	}

	result.PopulateCommon(resp)

	return result, err
}

func (client DatabasesClient) Get(ctx context.Context, databaseName string) (result Database, err error) {
	parameters := map[string]interface{}{
		"databaseName": databaseName,
	}

	req, err := client.Preparer(ctx, "Get", http.MethodGet, "/dbs/{databaseName}", parameters)
	if err != nil {
		return result, err
	}

	resp, err := client.Sender("Get", req)
	if err == nil {
		resp, err = client.Responder("Get", resp.Response, http.StatusOK, autorest.ByUnmarshallingJSON(&result))
		if err == nil {
			result.Path = "cosmos/dbs/" + databaseName
		}
	}
	result.Response = resp
	return result, err
}

func (client DatabasesClient) Delete(ctx context.Context, databaseName string) (result autorest.Response, err error) {
	parameters := map[string]interface{}{
		"databaseName": databaseName,
	}

	req, err := client.Preparer(ctx, "Delete", http.MethodDelete, "/dbs/{databaseName}", parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmos.DatabasesClient", "Delete", nil, "Failure preparing request")
		return result, err
	}

	resp, err := client.Sender("Delete", req)
	if err == nil {
		resp, err = client.Responder("Delete", resp.Response, http.StatusNoContent)
	}

	return resp, err
}
