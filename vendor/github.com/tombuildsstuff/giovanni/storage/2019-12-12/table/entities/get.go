package entities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type GetEntityInput struct {
	PartitionKey string
	RowKey       string

	// The Level of MetaData which should be returned
	MetaDataLevel MetaDataLevel
}

type GetEntityResult struct {
	autorest.Response

	Entity map[string]interface{}
}

// Get queries entities in a table and includes the $filter and $select options.
func (client Client) Get(ctx context.Context, accountName, tableName string, input GetEntityInput) (result GetEntityResult, err error) {
	if accountName == "" {
		return result, validation.NewError("entities.Client", "Get", "`accountName` cannot be an empty string.")
	}
	if tableName == "" {
		return result, validation.NewError("entities.Client", "Get", "`tableName` cannot be an empty string.")
	}
	if input.PartitionKey == "" {
		return result, validation.NewError("entities.Client", "Get", "`input.PartitionKey` cannot be an empty string.")
	}
	if input.RowKey == "" {
		return result, validation.NewError("entities.Client", "Get", "`input.RowKey` cannot be an empty string.")
	}

	req, err := client.GetPreparer(ctx, accountName, tableName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "entities.Client", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client Client) GetPreparer(ctx context.Context, accountName, tableName string, input GetEntityInput) (*http.Request, error) {

	pathParameters := map[string]interface{}{
		"tableName":    autorest.Encode("path", tableName),
		"partitionKey": autorest.Encode("path", input.PartitionKey),
		"rowKey":       autorest.Encode("path", input.RowKey),
	}

	headers := map[string]interface{}{
		"x-ms-version":          APIVersion,
		"Accept":                fmt.Sprintf("application/json;odata=%s", input.MetaDataLevel),
		"DataServiceVersion":    "3.0;NetFx",
		"MaxDataServiceVersion": "3.0;NetFx",
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetTableEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{tableName}(PartitionKey='{partitionKey}',RowKey='{rowKey}')", pathParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client Client) GetResponder(resp *http.Response) (result GetEntityResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Entity),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
