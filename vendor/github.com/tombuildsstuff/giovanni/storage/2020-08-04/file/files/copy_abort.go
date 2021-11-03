package files

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// AbortCopy aborts a pending Copy File operation, and leaves a destination file with zero length and full metadata
func (client Client) AbortCopy(ctx context.Context, accountName, shareName, path, fileName, copyID string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("files.Client", "AbortCopy", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("files.Client", "AbortCopy", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("files.Client", "AbortCopy", "`shareName` must be a lower-cased string.")
	}
	if fileName == "" {
		return result, validation.NewError("files.Client", "AbortCopy", "`fileName` cannot be an empty string.")
	}
	if copyID == "" {
		return result, validation.NewError("files.Client", "AbortCopy", "`copyID` cannot be an empty string.")
	}

	req, err := client.AbortCopyPreparer(ctx, accountName, shareName, path, fileName, copyID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "AbortCopy", nil, "Failure preparing request")
		return
	}

	resp, err := client.AbortCopySender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "files.Client", "AbortCopy", resp, "Failure sending request")
		return
	}

	result, err = client.AbortCopyResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "AbortCopy", resp, "Failure responding to request")
		return
	}

	return
}

// AbortCopyPreparer prepares the AbortCopy request.
func (client Client) AbortCopyPreparer(ctx context.Context, accountName, shareName, path, fileName, copyID string) (*http.Request, error) {
	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
		"directory": autorest.Encode("path", path),
		"fileName":  autorest.Encode("path", fileName),
	}

	queryParameters := map[string]interface{}{
		"comp":   autorest.Encode("query", "copy"),
		"copyid": autorest.Encode("query", copyID),
	}

	headers := map[string]interface{}{
		"x-ms-version":     APIVersion,
		"x-ms-copy-action": "abort",
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}{fileName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// AbortCopySender sends the AbortCopy request. The method will close the
// http.Response Body if it receives an error.
func (client Client) AbortCopySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// AbortCopyResponder handles the response to the AbortCopy request. The method always
// closes the http.Response Body.
func (client Client) AbortCopyResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
