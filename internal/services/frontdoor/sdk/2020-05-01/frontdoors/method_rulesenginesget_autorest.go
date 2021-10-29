package frontdoors

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type RulesEnginesGetResponse struct {
	HttpResponse *http.Response
	Model        *RulesEngine
}

// RulesEnginesGet ...
func (c FrontDoorsClient) RulesEnginesGet(ctx context.Context, id RulesEngineId) (result RulesEnginesGetResponse, err error) {
	req, err := c.preparerForRulesEnginesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRulesEnginesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRulesEnginesGet prepares the RulesEnginesGet request.
func (c FrontDoorsClient) preparerForRulesEnginesGet(ctx context.Context, id RulesEngineId) (*http.Request, error) {
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

// responderForRulesEnginesGet handles the response to the RulesEnginesGet request. The method always
// closes the http.Response Body.
func (c FrontDoorsClient) responderForRulesEnginesGet(resp *http.Response) (result RulesEnginesGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
