package shares

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// Delete deletes the specified Storage Share from within a Storage Account
func (client Client) Delete(ctx context.Context, accountName, shareName string, deleteSnapshots bool) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "Delete", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "Delete", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "Delete", "`shareName` must be a lower-cased string.")
	}

	req, err := client.DeletePreparer(ctx, accountName, shareName, deleteSnapshots)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "Delete", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client Client) DeletePreparer(ctx context.Context, accountName, shareName string, deleteSnapshots bool) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"restype": autorest.Encode("path", "share"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	if deleteSnapshots {
		headers["x-ms-delete-snapshots"] = "include"
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
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
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
