package forecast

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ExternalCloudProviderUsageResponse struct {
	HttpResponse *http.Response
	Model        *QueryResult
}

type ExternalCloudProviderUsageOptions struct {
	Filter *string
}

func DefaultExternalCloudProviderUsageOptions() ExternalCloudProviderUsageOptions {
	return ExternalCloudProviderUsageOptions{}
}

func (o ExternalCloudProviderUsageOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ExternalCloudProviderUsage ...
func (c ForecastClient) ExternalCloudProviderUsage(ctx context.Context, id ExternalCloudProviderTypeId, input ForecastDefinition, options ExternalCloudProviderUsageOptions) (result ExternalCloudProviderUsageResponse, err error) {
	req, err := c.preparerForExternalCloudProviderUsage(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "forecast.ForecastClient", "ExternalCloudProviderUsage", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "forecast.ForecastClient", "ExternalCloudProviderUsage", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForExternalCloudProviderUsage(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "forecast.ForecastClient", "ExternalCloudProviderUsage", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForExternalCloudProviderUsage prepares the ExternalCloudProviderUsage request.
func (c ForecastClient) preparerForExternalCloudProviderUsage(ctx context.Context, id ExternalCloudProviderTypeId, input ForecastDefinition, options ExternalCloudProviderUsageOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/forecast", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForExternalCloudProviderUsage handles the response to the ExternalCloudProviderUsage request. The method always
// closes the http.Response Body.
func (c ForecastClient) responderForExternalCloudProviderUsage(resp *http.Response) (result ExternalCloudProviderUsageResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
