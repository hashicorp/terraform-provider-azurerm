package fluxconfigurationoperationstatus

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type FluxConfigOperationStatusGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *OperationStatusResult
}

// FluxConfigOperationStatusGet ...
func (c FluxConfigurationOperationStatusClient) FluxConfigOperationStatusGet(ctx context.Context, id FluxConfigurationOperationId) (result FluxConfigOperationStatusGetOperationResponse, err error) {
	req, err := c.preparerForFluxConfigOperationStatusGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fluxconfigurationoperationstatus.FluxConfigurationOperationStatusClient", "FluxConfigOperationStatusGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "fluxconfigurationoperationstatus.FluxConfigurationOperationStatusClient", "FluxConfigOperationStatusGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForFluxConfigOperationStatusGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fluxconfigurationoperationstatus.FluxConfigurationOperationStatusClient", "FluxConfigOperationStatusGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForFluxConfigOperationStatusGet prepares the FluxConfigOperationStatusGet request.
func (c FluxConfigurationOperationStatusClient) preparerForFluxConfigOperationStatusGet(ctx context.Context, id FluxConfigurationOperationId) (*http.Request, error) {
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

// responderForFluxConfigOperationStatusGet handles the response to the FluxConfigOperationStatusGet request. The method always
// closes the http.Response Body.
func (c FluxConfigurationOperationStatusClient) responderForFluxConfigOperationStatusGet(resp *http.Response) (result FluxConfigOperationStatusGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
