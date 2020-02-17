package tables

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// Exists checks that the specified table exists
func (client Client) Exists(ctx context.Context, accountName, tableName string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("tables.Client", "Exists", "`accountName` cannot be an empty string.")
	}
	if tableName == "" {
		return result, validation.NewError("tables.Client", "Exists", "`tableName` cannot be an empty string.")
	}

	req, err := client.ExistsPreparer(ctx, accountName, tableName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tables.Client", "Exists", nil, "Failure preparing request")
		return
	}

	resp, err := client.ExistsSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "tables.Client", "Exists", resp, "Failure sending request")
		return
	}

	result, err = client.ExistsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tables.Client", "Exists", resp, "Failure responding to request")
		return
	}

	return
}

// ExistsPreparer prepares the Exists request.
func (client Client) ExistsPreparer(ctx context.Context, accountName, tableName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"tableName": autorest.Encode("path", tableName),
	}

	// NOTE: whilst the API documentation says that API Version is Optional
	// apparently specifying it causes an "invalid content type" to always be returned
	// as such we omit it here :shrug:

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.AsContentType("application/xml"),
		autorest.WithBaseURL(endpoints.GetTableEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/Tables('{tableName}')", pathParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ExistsSender sends the Exists request. The method will close the
// http.Response Body if it receives an error.
func (client Client) ExistsSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ExistsResponder handles the response to the Exists request. The method always
// closes the http.Response Body.
func (client Client) ExistsResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
