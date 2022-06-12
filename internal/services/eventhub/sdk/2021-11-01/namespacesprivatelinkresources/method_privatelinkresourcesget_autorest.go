package namespacesprivatelinkresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PrivateLinkResourcesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkResourcesListResult
}

// PrivateLinkResourcesGet ...
func (c NamespacesPrivateLinkResourcesClient) PrivateLinkResourcesGet(ctx context.Context, id NamespaceId) (result PrivateLinkResourcesGetOperationResponse, err error) {
	req, err := c.preparerForPrivateLinkResourcesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespacesprivatelinkresources.NamespacesPrivateLinkResourcesClient", "PrivateLinkResourcesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespacesprivatelinkresources.NamespacesPrivateLinkResourcesClient", "PrivateLinkResourcesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPrivateLinkResourcesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespacesprivatelinkresources.NamespacesPrivateLinkResourcesClient", "PrivateLinkResourcesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPrivateLinkResourcesGet prepares the PrivateLinkResourcesGet request.
func (c NamespacesPrivateLinkResourcesClient) preparerForPrivateLinkResourcesGet(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateLinkResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPrivateLinkResourcesGet handles the response to the PrivateLinkResourcesGet request. The method always
// closes the http.Response Body.
func (c NamespacesPrivateLinkResourcesClient) responderForPrivateLinkResourcesGet(resp *http.Response) (result PrivateLinkResourcesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
