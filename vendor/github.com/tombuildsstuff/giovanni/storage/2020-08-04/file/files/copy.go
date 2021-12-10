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

type CopyInput struct {
	// Specifies the URL of the source file or blob, up to 2 KB in length.
	//
	// To copy a file to another file within the same storage account, you may use Shared Key to authenticate
	// the source file. If you are copying a file from another storage account, or if you are copying a blob from
	// the same storage account or another storage account, then you must authenticate the source file or blob using a
	// shared access signature. If the source is a public blob, no authentication is required to perform the copy
	// operation. A file in a share snapshot can also be specified as a copy source.
	CopySource string

	MetaData map[string]string
}

type CopyResult struct {
	autorest.Response

	// The CopyID, which can be passed to AbortCopy to abort the copy.
	CopyID string

	// Either `success` or `pending`
	CopySuccess string
}

// Copy copies a blob or file to a destination file within the storage account asynchronously.
func (client Client) Copy(ctx context.Context, accountName, shareName, path, fileName string, input CopyInput) (result CopyResult, err error) {
	if accountName == "" {
		return result, validation.NewError("files.Client", "Copy", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("files.Client", "Copy", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("files.Client", "Copy", "`shareName` must be a lower-cased string.")
	}
	if fileName == "" {
		return result, validation.NewError("files.Client", "Copy", "`fileName` cannot be an empty string.")
	}
	if input.CopySource == "" {
		return result, validation.NewError("files.Client", "Copy", "`input.CopySource` cannot be an empty string.")
	}
	if err := metadata.Validate(input.MetaData); err != nil {
		return result, validation.NewError("files.Client", "Copy", fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
	}

	req, err := client.CopyPreparer(ctx, accountName, shareName, path, fileName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "Copy", nil, "Failure preparing request")
		return
	}

	resp, err := client.CopySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "files.Client", "Copy", resp, "Failure sending request")
		return
	}

	result, err = client.CopyResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "Copy", resp, "Failure responding to request")
		return
	}

	return
}

// CopyPreparer prepares the Copy request.
func (client Client) CopyPreparer(ctx context.Context, accountName, shareName, path, fileName string, input CopyInput) (*http.Request, error) {
	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
		"directory": autorest.Encode("path", path),
		"fileName":  autorest.Encode("path", fileName),
	}

	headers := map[string]interface{}{
		"x-ms-version":     APIVersion,
		"x-ms-copy-source": input.CopySource,
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

// CopySender sends the Copy request. The method will close the
// http.Response Body if it receives an error.
func (client Client) CopySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CopyResponder handles the response to the Copy request. The method always
// closes the http.Response Body.
func (client Client) CopyResponder(resp *http.Response) (result CopyResult, err error) {
	if resp != nil && resp.Header != nil {
		result.CopyID = resp.Header.Get("x-ms-copy-id")
		result.CopySuccess = resp.Header.Get("x-ms-copy-status")
	}

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
