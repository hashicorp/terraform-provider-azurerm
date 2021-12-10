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

type ClearByteRangeInput struct {
	StartBytes int64
	EndBytes   int64
}

// ClearByteRange clears the specified Byte Range from within the specified File
func (client Client) ClearByteRange(ctx context.Context, accountName, shareName, path, fileName string, input ClearByteRangeInput) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("files.Client", "ClearByteRange", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("files.Client", "ClearByteRange", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("files.Client", "ClearByteRange", "`shareName` must be a lower-cased string.")
	}
	if fileName == "" {
		return result, validation.NewError("files.Client", "ClearByteRange", "`fileName` cannot be an empty string.")
	}
	if input.StartBytes < 0 {
		return result, validation.NewError("files.Client", "ClearByteRange", "`input.StartBytes` must be greater or equal to 0.")
	}
	if input.EndBytes <= 0 {
		return result, validation.NewError("files.Client", "ClearByteRange", "`input.EndBytes` must be greater than 0.")
	}

	req, err := client.ClearByteRangePreparer(ctx, accountName, shareName, path, fileName, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "ClearByteRange", nil, "Failure preparing request")
		return
	}

	resp, err := client.ClearByteRangeSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "files.Client", "ClearByteRange", resp, "Failure sending request")
		return
	}

	result, err = client.ClearByteRangeResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "ClearByteRange", resp, "Failure responding to request")
		return
	}

	return
}

// ClearByteRangePreparer prepares the ClearByteRange request.
func (client Client) ClearByteRangePreparer(ctx context.Context, accountName, shareName, path, fileName string, input ClearByteRangeInput) (*http.Request, error) {
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
		"x-ms-version": APIVersion,
		"x-ms-write":   "clear",
		"x-ms-range":   fmt.Sprintf("bytes=%d-%d", input.StartBytes, input.EndBytes),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}{fileName}", pathParameters),
		autorest.WithHeaders(headers),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ClearByteRangeSender sends the ClearByteRange request. The method will close the
// http.Response Body if it receives an error.
func (client Client) ClearByteRangeSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ClearByteRangeResponder handles the response to the ClearByteRange request. The method always
// closes the http.Response Body.
func (client Client) ClearByteRangeResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
