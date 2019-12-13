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
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type CreateInput struct {
	// This header specifies the maximum size for the file, up to 1 TiB.
	ContentLength int64

	// The MIME content type of the file
	// If not specified, the default type is application/octet-stream.
	ContentType *string

	// Specifies which content encodings have been applied to the file.
	// This value is returned to the client when the Get File operation is performed
	// on the file resource and can be used to decode file content.
	ContentEncoding *string

	// Specifies the natural languages used by this resource.
	ContentLanguage *string

	// The File service stores this value but does not use or modify it.
	CacheControl *string

	// Sets the file's MD5 hash.
	ContentMD5 *string

	// Sets the file’s Content-Disposition header.
	ContentDisposition *string

	MetaData map[string]string
}

// Create creates a new file or replaces a file.
func (client Client) Create(ctx context.Context, accountName, shareName, path, fileName string, input CreateInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("files.Client", "Create", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("files.Client", "Create", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("files.Client", "Create", "`shareName` must be a lower-cased string.")
	}
	if fileName == "" {
		return result, validation.NewError("files.Client", "Create", "`fileName` cannot be an empty string.")
	}
	if err := metadata.Validate(input.MetaData); err != nil {
		return result, validation.NewError("files.Client", "Create", "`input.MetaData` cannot be an empty string.")
	}

	req, err := client.CreatePreparer(ctx, accountName, shareName, path, fileName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "files.Client", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "Create", resp, "Failure responding to request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client Client) CreatePreparer(ctx context.Context, accountName, shareName, path, fileName string, input CreateInput) (*http.Request, error) {
	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
		"directory": autorest.Encode("path", path),
		"fileName":  autorest.Encode("path", fileName),
	}

	headers := map[string]interface{}{
		"x-ms-version":        APIVersion,
		"x-ms-content-length": input.ContentLength,
		"x-ms-type":           "file",
	}

	if input.ContentDisposition != nil {
		headers["x-ms-content-disposition"] = *input.ContentDisposition
	}

	if input.ContentEncoding != nil {
		headers["x-ms-content-encoding"] = *input.ContentEncoding
	}

	if input.ContentMD5 != nil {
		headers["x-ms-content-md5"] = *input.ContentMD5
	}

	if input.ContentType != nil {
		headers["x-ms-content-type"] = *input.ContentType
	}

	headers = metadata.SetIntoHeaders(headers, input.MetaData)

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}{fileName}", pathParameters),
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
