package paths

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type PathResource string

const PathResourceFile PathResource = "file"
const PathResourceDirectory PathResource = "directory"

type CreateInput struct {
	Resource PathResource
}

// Create creates a Data Lake Store Gen2 Path within a Storage Account
func (client Client) Create(ctx context.Context, accountName string, fileSystemName string, path string, input CreateInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("datalakestore.Client", "Create", "`accountName` cannot be an empty string.")
	}
	if fileSystemName == "" {
		return result, validation.NewError("datalakestore.Client", "Create", "`fileSystemName` cannot be an empty string.")
	}

	req, err := client.CreatePreparer(ctx, accountName, fileSystemName, path, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "Create", resp, "Failure responding to request")
	}

	return
}

// CreatePreparer prepares the Create request.
func (client Client) CreatePreparer(ctx context.Context, accountName string, fileSystemName string, path string, input CreateInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fileSystemName": autorest.Encode("path", fileSystemName),
		"path":           autorest.Encode("path", path),
	}

	queryParameters := map[string]interface{}{
		"resource": autorest.Encode("query", input.Resource),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetDataLakeStoreEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{fileSystemName}/{path}", pathParameters),
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
