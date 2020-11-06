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

// DeleteSnapshot deletes the specified Snapshot of a Storage Share
func (client Client) DeleteSnapshot(ctx context.Context, accountName, shareName string, shareSnapshot string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "DeleteSnapshot", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "DeleteSnapshot", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "DeleteSnapshot", "`shareName` must be a lower-cased string.")
	}
	if shareSnapshot == "" {
		return result, validation.NewError("shares.Client", "DeleteSnapshot", "`shareSnapshot` cannot be an empty string.")
	}

	req, err := client.DeleteSnapshotPreparer(ctx, accountName, shareName, shareSnapshot)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "DeleteSnapshot", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSnapshotSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "DeleteSnapshot", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteSnapshotResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "DeleteSnapshot", resp, "Failure responding to request")
		return
	}

	return
}

// DeleteSnapshotPreparer prepares the DeleteSnapshot request.
func (client Client) DeleteSnapshotPreparer(ctx context.Context, accountName, shareName string, shareSnapshot string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"restype":       autorest.Encode("path", "share"),
		"sharesnapshot": autorest.Encode("query", shareSnapshot),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
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

// DeleteSnapshotSender sends the DeleteSnapshot request. The method will close the
// http.Response Body if it receives an error.
func (client Client) DeleteSnapshotSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteSnapshotResponder handles the response to the DeleteSnapshot request. The method always
// closes the http.Response Body.
func (client Client) DeleteSnapshotResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
