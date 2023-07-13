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

type InsertEntityInput struct {
	// The level of MetaData provided for this Entity
	MetaDataLevel MetaDataLevel

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

// Insert inserts a new entity into a table.
func (client Client) Insert(ctx context.Context, accountName, tableName string, input InsertEntityInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("entities.Client", "Insert", "`accountName` cannot be an empty string.")
	}
	if tableName == "" {
		return result, validation.NewError("entities.Client", "Insert", "`tableName` cannot be an empty string.")
	}
	if input.PartitionKey == "" {
		return result, validation.NewError("entities.Client", "Insert", "`input.PartitionKey` cannot be an empty string.")
	}
	if input.RowKey == "" {
		return result, validation.NewError("entities.Client", "Insert", "`input.RowKey` cannot be an empty string.")
	}

	req, err := client.InsertPreparer(ctx, accountName, tableName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "Insert", nil, "Failure preparing request")
		return
	}

	resp, err := client.InsertSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "entities.Client", "Insert", resp, "Failure sending request")
		return
	}

	result, err = client.InsertResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "Insert", resp, "Failure responding to request")
		return
	}

	return
}

// InsertPreparer prepares the Insert request.
func (client Client) InsertPreparer(ctx context.Context, accountName, tableName string, input InsertEntityInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"tableName": autorest.Encode("path", tableName),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
		"Accept":       fmt.Sprintf("application/json;odata=%s", input.MetaDataLevel),
		"Prefer":       "return-no-content",
	}

	input.Entity["PartitionKey"] = input.PartitionKey
	input.Entity["RowKey"] = input.RowKey

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json"),
		autorest.AsPost(),
		autorest.WithBaseURL(endpoints.GetTableEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{tableName}", pathParameters),
		autorest.WithJSON(input.Entity),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// InsertSender sends the Insert request. The method will close the
// http.Response Body if it receives an error.
func (client Client) InsertSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// InsertResponder handles the response to the Insert request. The method always
// closes the http.Response Body.
func (client Client) InsertResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
