package forecast

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

type UsageOptions struct {
	Filter *string
}

func DefaultUsageOptions() UsageOptions {
	return UsageOptions{}
}

func (o UsageOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// Usage ...
func (c ForecastClient) Usage(ctx context.Context, id ScopeId, input ForecastDefinition, options UsageOptions) (result UsageResponse, err error) {
	req, err := c.preparerForUsage(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "forecast.ForecastClient", "Usage", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "forecast.ForecastClient", "Usage", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUsage(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "forecast.ForecastClient", "Usage", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUsage prepares the Usage request.
func (c ForecastClient) preparerForUsage(ctx context.Context, id ScopeId, input ForecastDefinition, options UsageOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CostManagement/forecast", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUsage handles the response to the Usage request. The method always
// closes the http.Response Body.
func (c ForecastClient) responderForUsage(resp *http.Response) (result UsageResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
