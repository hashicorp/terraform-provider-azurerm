package shares

import (
	"context"
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
)

type setAcl struct {
	SignedIdentifiers []SignedIdentifier `xml:"SignedIdentifier"`

	XMLName xml.Name `xml:"SignedIdentifiers"`
}

// SetACL sets the specified Access Control List on the specified Storage Share
func (client Client) SetACL(ctx context.Context, accountName, shareName string, acls []SignedIdentifier) (result autorest.Response, err error) {
	if accountName == "" {
		return result, validation.NewError("shares.Client", "SetACL", "`accountName` cannot be an empty string.")
	}
	if shareName == "" {
		return result, validation.NewError("shares.Client", "SetACL", "`shareName` cannot be an empty string.")
	}
	if strings.ToLower(shareName) != shareName {
		return result, validation.NewError("shares.Client", "SetACL", "`shareName` must be a lower-cased string.")
	}

	req, err := client.SetACLPreparer(ctx, accountName, shareName, acls)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "SetACL", nil, "Failure preparing request")
		return
	}

	resp, err := client.SetACLSender(req)
	if err != nil {
		result = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "shares.Client", "SetACL", resp, "Failure sending request")
		return
	}

	result, err = client.SetACLResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "shares.Client", "SetACL", resp, "Failure responding to request")
		return
	}

	return
}

// SetACLPreparer prepares the SetACL request.
func (client Client) SetACLPreparer(ctx context.Context, accountName, shareName string, acls []SignedIdentifier) (*http.Request, error) {
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

	input := setAcl{
		SignedIdentifiers: acls,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/xml; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(endpoints.GetFileEndpoint(client.BaseURI, accountName)),
		autorest.WithPathParameters("/{shareName}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		autorest.WithHeaders(headers),
		autorest.WithXML(&input))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// SetACLSender sends the SetACL request. The method will close the
// http.Response Body if it receives an error.
func (client Client) SetACLSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// SetACLResponder handles the response to the SetACL request. The method always
// closes the http.Response Body.
func (client Client) SetACLResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result = autorest.Response{Response: resp}

	return
}
