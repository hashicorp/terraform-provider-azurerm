package entities

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type DeleteEntityInput struct {
	// When inserting an entity into a table, you must specify values for the PartitionKey and RowKey system properties.
	// Together, these properties form the primary key and must be unique within the table.
	// Both the PartitionKey and RowKey values must be string values; each key value may be up to 64 KB in size.
	// If you are using an integer value for the key value, you should convert the integer to a fixed-width string,
	// because they are canonically sorted. For example, you should convert the value 1 to 0000001 to ensure proper sorting.
	RowKey       string
	PartitionKey string
}

// Delete deletes an existing entity in a table.
func (client Client) Delete(ctx context.Context, accountName, tableName string, input DeleteEntityInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("entities.Client", "Delete", "`accountName` cannot be an empty string.")
	}
	if tableName == "" {
		return result, validation.NewError("entities.Client", "Delete", "`tableName` cannot be an empty string.")
	}
	if input.PartitionKey == "" {
		return result, validation.NewError("entities.Client", "Delete", "`input.PartitionKey` cannot be an empty string.")
	}
	if input.RowKey == "" {
		return result, validation.NewError("entities.Client", "Delete", "`input.RowKey` cannot be an empty string.")
	}

	req, err := client.DeletePreparer(ctx, accountName, tableName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "entities.Client", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "entities.Client", "Delete", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client Client) DeletePreparer(ctx context.Context, accountName, tableName string, input DeleteEntityInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"tableName":    autorest.Encode("path", tableName),
		"partitionKey": autorest.Encode("path", input.PartitionKey),
		"rowKey":       autorest.Encode("path", input.RowKey),
	}

	headers := map[string]interface{}{
		// TODO: support for eTags
		"If-Match": "*",
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(endpoints.GetTableEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{tableName}(PartitionKey='{partitionKey}', RowKey='{rowKey}')", pathParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client Client) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client Client) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
