package query

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type UsageResponse struct {
	HttpResponse *http.Response
	Model        *QueryResult
}

// Usage ...
func (c QueryClient) Usage(ctx context.Context, id ScopeId, input QueryDefinition) (result UsageResponse, err error) {
	req, err := c.preparerForUsage(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "query.QueryClient", "Usage", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "query.QueryClient", "Usage", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUsage(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "query.QueryClient", "Usage", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUsage prepares the Usage request.
func (c QueryClient) preparerForUsage(ctx context.Context, id ScopeId, input QueryDefinition) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CostManagement/query", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUsage handles the response to the Usage request. The method always
// closes the http.Response Body.
func (c QueryClient) responderForUsage(resp *http.Response) (result UsageResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
