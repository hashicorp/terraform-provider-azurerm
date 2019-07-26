package shares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetPropertiesResult struct {
	autorest.Response

	MetaData   map[string]string
	ShareQuota int
}

// GetProperties returns the properties about the specified Storage Share
func (client Client) GetProperties(ctx context.Context, accountName, shareName string) (result GetPropertiesResult, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "GetProperties", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "GetProperties", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "GetProperties", "`shareName` must be a lower-cased string.")
	}

	req, err := client.GetPropertiesPreparer(ctx, accountName, shareName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "GetProperties", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetPropertiesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "GetProperties", resp, "Failure sending request")
		return
	}

	result, err = client.GetPropertiesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "GetProperties", resp, "Failure responding to request")
		return
	}

	return
}

// GetPropertiesPreparer prepares the GetProperties request.
func (client Client) GetPropertiesPreparer(ctx context.Context, accountName, shareName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"restype": autorest.Encode("query", "share"),
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

// GetPropertiesSender sends the GetProperties request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetPropertiesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetPropertiesResponder handles the response to the GetProperties request. The method always
// closes the http.Response Body.
func (client Client) GetPropertiesResponder(resp *http.Response) (result GetPropertiesResult, err error) {
	if resp.Header != nil {
		result.MetaData = metadata.ParseFromHeaders(resp.Header)

		quotaRaw := resp.Header.Get("x-ms-share-quota")
		if quotaRaw != "" {
			quota, e := strconv.Atoi(quotaRaw)
			if e != nil {
				return result, fmt.Errorf("Error converting %q to an integer: %s", quotaRaw, err)
			}
			result.ShareQuota = quota
		}
	}

	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
