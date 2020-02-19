package shares

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type GetStatsResult struct {
	autorest.Response

	// The approximate size of the data stored on the share.
	// Note that this value may not include all recently created or recently resized files.
	ShareUsageBytes int64 `xml:"ShareUsageBytes"`
}

// GetStats returns information about the specified Storage Share
func (client Client) GetStats(ctx context.Context, accountName, shareName string) (result GetStatsResult, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "GetStats", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "GetStats", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "GetStats", "`shareName` must be a lower-cased string.")
	}

	req, err := client.GetStatsPreparer(ctx, accountName, shareName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "GetStats", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetStatsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "GetStats", resp, "Failure sending request")
		return
	}

	result, err = client.GetStatsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "GetStats", resp, "Failure responding to request")
		return
	}

	return
}

// GetStatsPreparer prepares the GetStats request.
func (client Client) GetStatsPreparer(ctx context.Context, accountName, shareName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"restype": autorest.Encode("query", "share"),
		"comp":    autorest.Encode("query", "stats"),
	}

	headers := map[string]interface{}{
		"x-ms-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetStatsSender sends the GetStats request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetStatsSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetStatsResponder handles the response to the GetStats request. The method always
// closes the http.Response Body.
func (client Client) GetStatsResponder(resp *http.Response) (result GetStatsResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingXML(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
