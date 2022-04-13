package shares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type CreateSnapshotInput struct {
	MetaData map[string]string
}

type CreateSnapshotResult struct {
	autorest.Response

	// This header is a DateTime value that uniquely identifies the share snapshot.
	// The value of this header may be used in subsequent requests to access the share snapshot.
	// This value is opaque.
	SnapshotDateTime string
}

// CreateSnapshot creates a read-only snapshot of the share
// A share can support creation of 200 share snapshots. Attempting to create more than 200 share snapshots fails with 409 (Conflict).
// Attempting to create a share snapshot while a previous Snapshot Share operation is in progress fails with 409 (Conflict).
func (client Client) CreateSnapshot(ctx context.Context, accountName, shareName string, input CreateSnapshotInput) (result CreateSnapshotResult, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "CreateSnapshot", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "CreateSnapshot", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "CreateSnapshot", "`shareName` must be a lower-cased string.")
	}
	if err := metadata.Validate(input.MetaData); err != nil {
		return result, validation.NewError("shares.Client", "CreateSnapshot", fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
	}

	req, err := client.CreateSnapshotPreparer(ctx, accountName, shareName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "CreateSnapshot", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSnapshotSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "CreateSnapshot", resp, "Failure sending request")
		return
	}

	result, err = client.CreateSnapshotResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "CreateSnapshot", resp, "Failure responding to request")
		return
	}

	return
}

// CreateSnapshotPreparer prepares the CreateSnapshot request.
func (client Client) CreateSnapshotPreparer(ctx context.Context, accountName, shareName string, input CreateSnapshotInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"comp":    autorest.Encode("query", "snapshot"),
		"restype": autorest.Encode("query", "share"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	headers = metadata.SetIntoHeaders(headers, input.MetaData)

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSnapshotSender sends the CreateSnapshot request. The method will close the
// http.Response Body if it receives an error.
func (client Client) CreateSnapshotSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateSnapshotResponder handles the response to the CreateSnapshot request. The method always
// closes the http.Response Body.
func (client Client) CreateSnapshotResponder(resp *http.Response) (result CreateSnapshotResult, err error) {
	result.SnapshotDateTime = resp.Header.Get("x-ms-snapshot")

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
