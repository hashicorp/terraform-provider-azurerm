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

type GetByteRangeInput struct {
	StartBytes int64
	EndBytes   int64
}

type GetByteRangeResult struct {
	autorest.Response

	Contents []byte
}

// GetByteRange returns the specified Byte Range from the specified File.
func (client Client) GetByteRange(ctx context.Context, accountName, shareName, path, fileName string, input GetByteRangeInput) (result GetByteRangeResult, err error) {
	if accountName == "" {
		return result, validation.NewError("files.Client", "GetByteRange", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("files.Client", "GetByteRange", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("files.Client", "GetByteRange", "`shareName` must be a lower-cased string.")
	}
	if fileName == "" {
		return result, validation.NewError("files.Client", "GetByteRange", "`fileName` cannot be an empty string.")
	}
	if input.StartBytes < 0 {
		return result, validation.NewError("files.Client", "GetByteRange", "`input.StartBytes` must be greater or equal to 0.")
	}
	if input.EndBytes <= 0 {
		return result, validation.NewError("files.Client", "GetByteRange", "`input.EndBytes` must be greater than 0.")
	}
	expectedBytes := input.EndBytes - input.StartBytes
	if expectedBytes < (4 * 1024) {
		return result, validation.NewError("files.Client", "GetByteRange", "Requested Byte Range must be at least 4KB.")
	}
	if expectedBytes > (4 * 1024 * 1024) {
		return result, validation.NewError("files.Client", "GetByteRange", "Requested Byte Range must be at most 4MB.")
	}

	req, err := client.GetByteRangePreparer(ctx, accountName, shareName, path, fileName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "GetByteRange", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetByteRangeSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "files.Client", "GetByteRange", resp, "Failure sending request")
		return
	}

	result, err = client.GetByteRangeResponder(resp, expectedBytes)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "GetByteRange", resp, "Failure responding to request")
		return
	}

	return
}

// GetByteRangePreparer prepares the GetByteRange request.
func (client Client) GetByteRangePreparer(ctx context.Context, accountName, shareName, path, fileName string, input GetByteRangeInput) (*http.Request, error) {
	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
		"directory": autorest.Encode("path", path),
		"fileName":  autorest.Encode("path", fileName),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
		"x-ms-range":   fmt.Sprintf("bytes=%d-%d", input.StartBytes, input.EndBytes-1),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}{fileName}", pathParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetByteRangeSender sends the GetByteRange request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetByteRangeSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetByteRangeResponder handles the response to the GetByteRange request. The method always
// closes the http.Response Body.
func (client Client) GetByteRangeResponder(resp *http.Response, length int64) (result GetByteRangeResult, err error) {
	result.Contents = make([]byte, length)

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusPartialContent),
		autorest.ByUnmarshallingBytes(&result.Contents),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
