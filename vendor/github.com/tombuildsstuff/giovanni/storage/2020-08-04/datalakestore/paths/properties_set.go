package paths

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type SetAccessControlInput struct {
	Owner *string
	Group *string
	ACL   *string

	// Optional - A date and time value.
	// Specify this header to perform the operation only if the resource has been modified since the specified date and time.
	IfModifiedSince *string

	// Optional - A date and time value.
	// Specify this header to perform the operation only if the resource has not been modified since the specified date and time.
	IfUnmodifiedSince *string
}

// SetProperties sets the access control properties for a Data Lake Store Gen2 Path within a Storage Account File System
func (client Client) SetAccessControl(ctx context.Context, accountName string, fileSystemName string, path string, input SetAccessControlInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("datalakestore.Client", "SetAccessControl", "`accountName` cannot be an empty string.")
	}
	if fileSystemName == "" {
		return result, validation.NewError("datalakestore.Client", "SetAccessControl", "`fileSystemName` cannot be an empty string.")
	}

	req, err := client.SetAccessControlPreparer(ctx, accountName, fileSystemName, path, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "SetAccessControl", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetAccessControlSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "SetAccessControl", resp, "Failure sending request")
		return
	}

	result, err = client.SetAccessControlResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datalakestore.Client", "SetAccessControl", resp, "Failure responding to request")
	}

	return
}

// SetAccessControlPreparer prepares the SetAccessControl request.
func (client Client) SetAccessControlPreparer(ctx context.Context, accountName string, fileSystemName string, path string, input SetAccessControlInput) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"fileSystemName": autorest.Encode("path", fileSystemName),
		"path":           autorest.Encode("path", path),
	}

	queryParameters := map[string]interface{}{
		"action": autorest.Encode("query", "setAccessControl"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	if input.Owner != nil {
		headers["x-ms-owner"] = *input.Owner
	}
	if input.Group != nil {
		headers["x-ms-group"] = *input.Group
	}
	if input.ACL != nil {
		headers["x-ms-acl"] = *input.ACL
	}

	if input.IfModifiedSince != nil {
		headers["If-Modified-Since"] = *input.IfModifiedSince
	}
	if input.IfUnmodifiedSince != nil {
		headers["If-Unmodified-Since"] = *input.IfUnmodifiedSince
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPatch(),
		autorest.WithBaseURL(endpoints.GetDataLakeStoreEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{fileSystemName}/{path}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))

	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetAccessControlSender sends the SetAccessControl request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetAccessControlSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetAccessControlResponder handles the response to the SetAccessControl request. The method always
// closes the http.Response Body.
func (client Client) SetAccessControlResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}
	return
}
