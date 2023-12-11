package paths

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

// Delete deletes a Data Lake Store Gen2 FileSystem within a Storage Account
func (client Client) Delete(ctx context.Context, accountName string, fileSystemName string, path string) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("datalakestore.Client", "Delete", "`accountName` cannot be an empty string.")
	}
	if fileSystemName == "" {
		return result, validation.NewError("datalakestore.Client", "Delete", "`fileSystemName` cannot be an empty string.")
	}

	req, err := client.DeletePreparer(ctx, accountName, fileSystemName, path)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client Client) DeletePreparer(ctx context.Context, accountName string, fileSystemName string, path string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fileSystemName": autorest.Encode("path", fileSystemName),
		"path":           autorest.Encode("path", path),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(endpoints.GetDataLakeStoreEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{fileSystemName}/{path}", pathParameters),
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
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
