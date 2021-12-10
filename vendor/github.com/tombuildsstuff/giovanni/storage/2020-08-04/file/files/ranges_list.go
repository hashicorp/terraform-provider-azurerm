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

type ListRangesResult struct {
	autorest.Response

	Ranges []Range `xml:"Range"`
}

type Range struct {
	Start string `xml:"Start"`
	End   string `xml:"End"`
}

// ListRanges returns the list of valid ranges for the specified File.
func (client Client) ListRanges(ctx context.Context, accountName, shareName, path, fileName string) (result ListRangesResult, err error) {
	if accountName == "" {
		return result, validation.NewError("files.Client", "ListRanges", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("files.Client", "ListRanges", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("files.Client", "ListRanges", "`shareName` must be a lower-cased string.")
	}
	if path == "" {
		return result, validation.NewError("files.Client", "ListRanges", "`path` cannot be an empty string.")
	}
	if fileName == "" {
		return result, validation.NewError("files.Client", "ListRanges", "`fileName` cannot be an empty string.")
	}

	req, err := client.ListRangesPreparer(ctx, accountName, shareName, path, fileName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "ListRanges", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListRangesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "files.Client", "ListRanges", resp, "Failure sending request")
		return
	}

	result, err = client.ListRangesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "files.Client", "ListRanges", resp, "Failure responding to request")
		return
	}

	return
}

// ListRangesPreparer prepares the ListRanges request.
func (client Client) ListRangesPreparer(ctx context.Context, accountName, shareName, path, fileName string) (*http.Request, error) {
	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
		"directory": autorest.Encode("path", path),
		"fileName":  autorest.Encode("path", fileName),
	}

	queryParameters := map[string]interface{}{
		"comp": autorest.Encode("query", "rangelist"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}/{directory}{fileName}", pathParameters),
		autorest.WithHeaders(headers),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListRangesSender sends the ListRanges request. The method will close the
// http.Response Body if it receives an error.
func (client Client) ListRangesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListRangesResponder handles the response to the ListRanges request. The method always
// closes the http.Response Body.
func (client Client) ListRangesResponder(resp *http.Response) (result ListRangesResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingXML(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
