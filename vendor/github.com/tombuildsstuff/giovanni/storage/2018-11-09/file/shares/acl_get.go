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

type GetACLResult struct {
	autorest.Response

	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`
}

// GetACL get the Access Control List for the specified Storage Share
func (client Client) GetACL(ctx context.Context, accountName, shareName string) (result GetACLResult, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "GetACL", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "GetACL", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "GetACL", "`shareName` must be a lower-cased string.")
	}

	req, err := client.GetACLPreparer(ctx, accountName, shareName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "GetACL", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetACLSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "GetACL", resp, "Failure sending request")
		return
	}

	result, err = client.GetACLResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "GetACL", resp, "Failure responding to request")
		return
	}

	return
}

// GetACLPreparer prepares the GetACL request.
func (client Client) GetACLPreparer(ctx context.Context, accountName, shareName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"shareName": autorest.Encode("path", shareName),
	}

	queryParameters := map[string]interface{}{
		"restype": autorest.Encode("query", "share"),
		"comp":    autorest.Encode("query", "acl"),
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

// GetACLSender sends the GetACL request. The method will close the
// http.Response Body if it receives an error.
func (client Client) GetACLSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetACLResponder handles the response to the GetACL request. The method always
// closes the http.Response Body.
func (client Client) GetACLResponder(resp *http.Response) (result GetACLResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingXML(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}

	return
}
