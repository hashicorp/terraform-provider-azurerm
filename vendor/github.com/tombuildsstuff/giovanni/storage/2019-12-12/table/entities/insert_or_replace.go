package entities

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type InsertOrReplaceEntityInput struct {
	// The Entity which should be inserted, by default all values are strings
	// To explicitly type a property, specify the appropriate OData data type by setting
	// the m:type attribute within the property definition
	Entity map[string]interface{}

	// When inserting an entity into a table, you must specify values for the PartitionKey and RowKey system properties.
	// Together, these properties form the primary key and must be unique within the table.
	// Both the PartitionKey and RowKey values must be string values; each key value may be up to 64 KB in size.
	// If you are using an integer value for the key value, you should convert the integer to a fixed-width string,
	// because they are canonically sorted. For example, you should convert the value 1 to 0000001 to ensure proper sorting.
	RowKey       string
	PartitionKey string
}

// InsertOrReplace replaces an existing entity or inserts a new entity if it does not exist in the table.
// Because this operation can insert or update an entity, it is also known as an upsert operation.
func (client Client) InsertOrReplace(ctx context.Context, accountName, tableName string, input InsertOrReplaceEntityInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("entities.Client", "InsertOrReplace", "`accountName` cannot be an empty string.")
	}
	if tableName == "" {
		return result, validation.NewError("entities.Client", "InsertOrReplace", "`tableName` cannot be an empty string.")
	}
	if input.PartitionKey == "" {
		return result, validation.NewError("entities.Client", "InsertOrReplace", "`input.PartitionKey` cannot be an empty string.")
	}
	if input.RowKey == "" {
		return result, validation.NewError("entities.Client", "InsertOrReplace", "`input.RowKey` cannot be an empty string.")
	}

	req, err := client.InsertOrReplacePreparer(ctx, accountName, tableName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "InsertOrReplace", nil, "Failure preparing request")
		return
	}

	resp, err := client.InsertOrReplaceSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "entities.Client", "InsertOrReplace", resp, "Failure sending request")
		return
	}

	result, err = client.InsertOrReplaceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "InsertOrReplace", resp, "Failure responding to request")
		return
	}

	return
}

// InsertOrReplacePreparer prepares the InsertOrReplace request.
func (client Client) InsertOrReplacePreparer(ctx context.Context, accountName, tableName string, input InsertOrReplaceEntityInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"tableName":    autorest.Encode("path", tableName),
		"partitionKey": autorest.Encode("path", input.PartitionKey),
		"rowKey":       autorest.Encode("path", input.RowKey),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
		"Accept":       "application/json",
		"Prefer":       "return-no-content",
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json"),
		autorest.AsMerge(),
		autorest.WithBaseURL(endpoints.GetTableEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{tableName}(PartitionKey='{partitionKey}', RowKey='{rowKey}')", pathParameters),
		autorest.WithJSON(input.Entity),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// InsertOrReplaceSender sends the InsertOrReplace request. The method will close the
// http.Response Body if it receives an error.
func (client Client) InsertOrReplaceSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// InsertOrReplaceResponder handles the response to the InsertOrReplace request. The method always
// closes the http.Response Body.
func (client Client) InsertOrReplaceResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
