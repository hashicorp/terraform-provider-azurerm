package extensionoperationstatus

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type OperationStatusGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *OperationStatusResult
}

// OperationStatusGet ...
func (c ExtensionOperationStatusClient) OperationStatusGet(ctx context.Context, id OperationId) (result OperationStatusGetOperationResponse, err error) {
	req, err := c.preparerForOperationStatusGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensionoperationstatus.ExtensionOperationStatusClient", "OperationStatusGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensionoperationstatus.ExtensionOperationStatusClient", "OperationStatusGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForOperationStatusGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensionoperationstatus.ExtensionOperationStatusClient", "OperationStatusGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForOperationStatusGet prepares the OperationStatusGet request.
func (c ExtensionOperationStatusClient) preparerForOperationStatusGet(ctx context.Context, id OperationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForOperationStatusGet handles the response to the OperationStatusGet request. The method always
// closes the http.Response Body.
func (c ExtensionOperationStatusClient) responderForOperationStatusGet(resp *http.Response) (result OperationStatusGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
