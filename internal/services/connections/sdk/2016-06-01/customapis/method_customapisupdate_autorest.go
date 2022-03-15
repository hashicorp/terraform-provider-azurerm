package customapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CustomApisUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *CustomApiDefinition
}

// CustomApisUpdate ...
func (c CustomAPIsClient) CustomApisUpdate(ctx context.Context, id CustomApiId, input CustomApiDefinition) (result CustomApisUpdateOperationResponse, err error) {
	req, err := c.preparerForCustomApisUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomApisUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customapis.CustomAPIsClient", "CustomApisUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomApisUpdate prepares the CustomApisUpdate request.
func (c CustomAPIsClient) preparerForCustomApisUpdate(ctx context.Context, id CustomApiId, input CustomApiDefinition) (*http.Request, error) {
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

// responderForCustomApisUpdate handles the response to the CustomApisUpdate request. The method always
// closes the http.Response Body.
func (c CustomAPIsClient) responderForCustomApisUpdate(resp *http.Response) (result CustomApisUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
