package namespaces

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PatchOperationResponse struct {
	HttpResponse *http.Response
	Model        *NamespaceResource
}

// Patch ...
func (c NamespacesClient) Patch(ctx context.Context, id NamespaceId, input NamespacePatchParameters) (result PatchOperationResponse, err error) {
	req, err := c.preparerForPatch(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "Patch", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "Patch", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPatch(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "Patch", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPatch prepares the Patch request.
func (c NamespacesClient) preparerForPatch(ctx context.Context, id NamespaceId, input NamespacePatchParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPatch handles the response to the Patch request. The method always
// closes the http.Response Body.
func (c NamespacesClient) responderForPatch(resp *http.Response) (result PatchOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
