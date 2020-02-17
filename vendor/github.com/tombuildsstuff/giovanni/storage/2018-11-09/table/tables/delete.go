package tables

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// Delete deletes the specified table and any data it contains.
func (client Client) Delete(ctx context.Context, accountName, tableName string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("tables.Client", "Delete", "`accountName` cannot be an empty string.")
	}
	if tableName == "" {
		return result, validation.NewError("tables.Client", "Delete", "`tableName` cannot be an empty string.")
	}

	req, err := client.DeletePreparer(ctx, accountName, tableName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tables.Client", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "tables.Client", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tables.Client", "Delete", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client Client) DeletePreparer(ctx context.Context, accountName, tableName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"tableName": autorest.Encode("path", tableName),
	}

	// NOTE: whilst the API documentation says that API Version is Optional
	// apparently specifying it causes an "invalid content type" to always be returned
	// as such we omit it here :shrug:

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(endpoints.GetTableEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/Tables('{tableName}')", pathParameters))
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
