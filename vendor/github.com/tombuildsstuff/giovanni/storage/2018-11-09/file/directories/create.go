package directories

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

// Create creates a new directory under the specified share or parent directory.
func (client Client) Create(ctx context.Context, accountName, shareName, path string, metaData map[string]string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("directories.Client", "Create", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("directories.Client", "Create", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("directories.Client", "Create", "`shareName` must be a lower-cased string.")
	}
	if path == "" {
		return result, validation.NewError("directories.Client", "Create", "`path` cannot be an empty string.")
	}
	if err := metadata.Validate(metaData); err != nil {
		return result, validation.NewError("directories.Client", "Create", fmt.Sprintf("`metadata` is not valid: %s.", err))
	}

	req, err := client.CreatePreparer(ctx, accountName, shareName, path, metaData)
	if err != nil {
		err = autorest.NewErrorWithError(err, "directories.Client", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "directories.Client", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "directories.Client", "Create", resp, "Failure responding to request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client Client) CreatePreparer(ctx context.Context, accountName, shareName, path string, metaData map[string]string) (*http.Request, error) {
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

	headers = metadata.SetIntoHeaders(headers, metaData)

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client Client) CreateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client Client) CreateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
