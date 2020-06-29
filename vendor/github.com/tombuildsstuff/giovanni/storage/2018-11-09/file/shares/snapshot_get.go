package shares

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetSnapshotPropertiesResult struct {
	autorest.Response

	MetaData map[string]string
}

// GetSnapshot gets information about the specified Snapshot of the specified Storage Share
func (client Client) GetSnapshot(ctx context.Context, accountName, shareName, snapshotShare string) (result GetSnapshotPropertiesResult, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "GetSnapshot", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "GetSnapshot", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "GetSnapshot", "`shareName` must be a lower-cased string.")
	}
	if snapshotShare == "" {
		return result, validation.NewError("shares.Client", "GetSnapshot", "`snapshotShare` cannot be an empty string.")
	}

	req, err := client.GetSnapshotPreparer(ctx, accountName, shareName, snapshotShare)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "GetSnapshot", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSnapshotSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "GetSnapshot", resp, "Failure sending request")
		return
	}

	result, err = client.GetSnapshotResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "GetSnapshot", resp, "Failure responding to request")
		return
	}

	return
}

// GetSnapshotPreparer prepares the GetSnapshot request.
func (client Client) GetSnapshotPreparer(ctx context.Context, accountName, shareName, snapshotShare string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"restype":  autorest.Encode("query", "share"),
		"snapshot": autorest.Encode("query", snapshotShare),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSnapshotSender sends the GetSnapshot request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetSnapshotSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetSnapshotResponder handles the response to the GetSnapshot request. The method always
// closes the http.Response Body.
func (client Client) GetSnapshotResponder(resp *http.Response) (result GetSnapshotPropertiesResult, err error) {
	if resp.Header != nil {
		result.MetaData = metadata.ParseFromHeaders(resp.Header)
	}

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
