package customapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CustomApisCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *CustomApiDefinition
}

// CustomApisCreateOrUpdate ...
func (c CustomAPIsClient) CustomApisCreateOrUpdate(ctx context.Context, id CustomApiId, input CustomApiDefinition) (result CustomApisCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForCustomApisCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomApisCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomApisCreateOrUpdate prepares the CustomApisCreateOrUpdate request.
func (c CustomAPIsClient) preparerForCustomApisCreateOrUpdate(ctx context.Context, id CustomApiId, input CustomApiDefinition) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCustomApisCreateOrUpdate handles the response to the CustomApisCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c CustomAPIsClient) responderForCustomApisCreateOrUpdate(resp *http.Response) (result CustomApisCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
