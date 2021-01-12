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

type PutByteRangeInput struct {
	StartBytes int64
	EndBytes   int64

	// Content is the File Contents for the specified range
	// which can be at most 4MB
	Content []byte
}

// PutByteRange puts the specified Byte Range in the specified File.
func (client Client) PutByteRange(ctx context.Context, accountName, shareName, path, fileName string, input PutByteRangeInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("files.Client", "PutByteRange", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("files.Client", "PutByteRange", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("files.Client", "PutByteRange", "`shareName` must be a lower-cased string.")
	}
	if fileName == "" {
		return result, validation.NewError("files.Client", "PutByteRange", "`fileName` cannot be an empty string.")
	}
	if input.StartBytes < 0 {
		return result, validation.NewError("files.Client", "PutByteRange", "`input.StartBytes` must be greater or equal to 0.")
	}
	if input.EndBytes <= 0 {
		return result, validation.NewError("files.Client", "PutByteRange", "`input.EndBytes` must be greater than 0.")
	}

	expectedBytes := input.EndBytes - input.StartBytes
	actualBytes := len(input.Content)
	if expectedBytes != int64(actualBytes) {
		return result, validation.NewError("files.Client", "PutByteRange", fmt.Sprintf("The specified byte-range (%d) didn't match the content size (%d).", expectedBytes, actualBytes))
	}

	if expectedBytes > (4 * 1024 * 1024) {
		return result, validation.NewError("files.Client", "PutByteRange", "Specified Byte Range must be at most 4MB.")
	}

	req, err := client.PutByteRangePreparer(ctx, accountName, shareName, path, fileName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "PutByteRange", nil, "Failure preparing request")
		return
	}

	resp, err := client.PutByteRangeSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "files.Client", "PutByteRange", resp, "Failure sending request")
		return
	}

	result, err = client.PutByteRangeResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "PutByteRange", resp, "Failure responding to request")
		return
	}

	return
}

// PutByteRangePreparer prepares the PutByteRange request.
func (client Client) PutByteRangePreparer(ctx context.Context, accountName, shareName, path, fileName string, input PutByteRangeInput) (*http.Request, error) {
	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
		"directory": autorest.Encode("path", path),
		"fileName":  autorest.Encode("path", fileName),
	}

	queryParameters := map[string]interface{}{
		"comp": autorest.Encode("query", "range"),
	}

	headers := map[string]interface{}{
		"x-ms-version":   APIVersion,
		"x-ms-write":     "update",
		"x-ms-range":     fmt.Sprintf("bytes=%d-%d", input.StartBytes, input.EndBytes-1),
		"Content-Length": int(len(input.Content)),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}{fileName}", pathParameters),
		autorest.WithHeaders(headers),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithBytes(&input.Content))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PutByteRangeSender sends the PutByteRange request. The method will close the
// http.Response Body if it receives an error.
func (client Client) PutByteRangeSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// PutByteRangeResponder handles the response to the PutByteRange request. The method always
// closes the http.Response Body.
func (client Client) PutByteRangeResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
