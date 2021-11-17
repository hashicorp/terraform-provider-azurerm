package partnernamespaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListSharedAccessKeysResponse struct {
	HttpResponse *http.Response
	Model        *PartnerNamespaceSharedAccessKeys
}

// ListSharedAccessKeys ...
func (c PartnerNamespacesClient) ListSharedAccessKeys(ctx context.Context, id PartnerNamespaceId) (result ListSharedAccessKeysResponse, err error) {
	req, err := c.preparerForListSharedAccessKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "partnernamespaces.PartnerNamespacesClient", "ListSharedAccessKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "partnernamespaces.PartnerNamespacesClient", "ListSharedAccessKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListSharedAccessKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "partnernamespaces.PartnerNamespacesClient", "ListSharedAccessKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListSharedAccessKeys prepares the ListSharedAccessKeys request.
func (c PartnerNamespacesClient) preparerForListSharedAccessKeys(ctx context.Context, id PartnerNamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListSharedAccessKeys handles the response to the ListSharedAccessKeys request. The method always
// closes the http.Response Body.
func (c PartnerNamespacesClient) responderForListSharedAccessKeys(resp *http.Response) (result ListSharedAccessKeysResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
