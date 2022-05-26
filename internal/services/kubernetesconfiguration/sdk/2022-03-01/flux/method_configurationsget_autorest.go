package flux

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConfigurationsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *FluxConfiguration
}

// ConfigurationsGet ...
func (c FluxClient) ConfigurationsGet(ctx context.Context, id FluxConfigurationId) (result ConfigurationsGetOperationResponse, err error) {
	req, err := c.preparerForConfigurationsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "flux.FluxClient", "ConfigurationsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "flux.FluxClient", "ConfigurationsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConfigurationsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "flux.FluxClient", "ConfigurationsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConfigurationsGet prepares the ConfigurationsGet request.
func (c FluxClient) preparerForConfigurationsGet(ctx context.Context, id FluxConfigurationId) (*http.Request, error) {
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

// responderForConfigurationsGet handles the response to the ConfigurationsGet request. The method always
// closes the http.Response Body.
func (c FluxClient) responderForConfigurationsGet(resp *http.Response) (result ConfigurationsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
