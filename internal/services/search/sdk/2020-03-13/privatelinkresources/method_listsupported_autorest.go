package privatelinkresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListSupportedOperationResponse struct {
	HttpResponse *http.Response
	Model        *PrivateLinkResourcesResult
}

type ListSupportedOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultListSupportedOperationOptions() ListSupportedOperationOptions {
	return ListSupportedOperationOptions{}
}

func (o ListSupportedOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.XMsClientRequestId != nil {
		out["x-ms-client-request-id"] = *o.XMsClientRequestId
	}

	return out
}

func (o ListSupportedOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// ListSupported ...
func (c PrivateLinkResourcesClient) ListSupported(ctx context.Context, id SearchServiceId, options ListSupportedOperationOptions) (result ListSupportedOperationResponse, err error) {
	req, err := c.preparerForListSupported(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkresources.PrivateLinkResourcesClient", "ListSupported", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkresources.PrivateLinkResourcesClient", "ListSupported", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListSupported(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkresources.PrivateLinkResourcesClient", "ListSupported", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListSupported prepares the ListSupported request.
func (c PrivateLinkResourcesClient) preparerForListSupported(ctx context.Context, id SearchServiceId, options ListSupportedOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/privateLinkResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListSupported handles the response to the ListSupported request. The method always
// closes the http.Response Body.
func (c PrivateLinkResourcesClient) responderForListSupported(resp *http.Response) (result ListSupportedOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
