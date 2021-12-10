package directories

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

type GetResult struct {
	autorest.Response

	// A set of name-value pairs that contain metadata for the directory.
	MetaData map[string]string

	// The value of this header is set to true if the directory metadata is completely
	// encrypted using the specified algorithm. Otherwise, the value is set to false.
	DirectoryMetaDataEncrypted bool
}

// Get returns all system properties for the specified directory,
// and can also be used to check the existence of a directory.
func (client Client) Get(ctx context.Context, accountName, shareName, path string) (result GetResult, err error) {
	if accountName == "" {
		return result, validation.NewError("directories.Client", "Get", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("directories.Client", "Get", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("directories.Client", "Get", "`shareName` must be a lower-cased string.")
	}
	if path == "" {
		return result, validation.NewError("directories.Client", "Get", "`path` cannot be an empty string.")
	}

	req, err := client.GetPreparer(ctx, accountName, shareName, path)
	if err != nil {
		err = autorest.NewErrorWithError(err, "directories.Client", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "directories.Client", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "directories.Client", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client Client) GetPreparer(ctx context.Context, accountName, shareName, path string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
		"directory": autorest.Encode("path", path),
	}

	queryParameters := map[string]interface{}{
		"restype": autorest.Encode("query", "directory"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client Client) GetResponder(resp *http.Response) (result GetResult, err error) {
	if resp != nil && resp.Header != nil {
		result.MetaData = metadata.ParseFromHeaders(resp.Header)
		result.DirectoryMetaDataEncrypted = strings.EqualFold(resp.Header.Get("x-ms-server-encrypted"), "true")
	}

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
