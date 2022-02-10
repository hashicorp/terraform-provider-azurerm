package heatmaps

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type HeatMapGetResponse struct {
	HttpResponse *http.Response
	Model        *HeatMapModel
}

type HeatMapGetOptions struct {
	BotRight *string
	TopLeft  *string
}

func DefaultHeatMapGetOptions() HeatMapGetOptions {
	return HeatMapGetOptions{}
}

func (o HeatMapGetOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.BotRight != nil {
		out["botRight"] = *o.BotRight
	}

	if o.TopLeft != nil {
		out["topLeft"] = *o.TopLeft
	}

	return out
}

// HeatMapGet ...
func (c HeatMapsClient) HeatMapGet(ctx context.Context, id HeatMapTypeId, options HeatMapGetOptions) (result HeatMapGetResponse, err error) {
	req, err := c.preparerForHeatMapGet(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "heatmaps.HeatMapsClient", "HeatMapGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "heatmaps.HeatMapsClient", "HeatMapGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForHeatMapGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "heatmaps.HeatMapsClient", "HeatMapGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForHeatMapGet prepares the HeatMapGet request.
func (c HeatMapsClient) preparerForHeatMapGet(ctx context.Context, id HeatMapTypeId, options HeatMapGetOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForHeatMapGet handles the response to the HeatMapGet request. The method always
// closes the http.Response Body.
func (c HeatMapsClient) responderForHeatMapGet(resp *http.Response) (result HeatMapGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
