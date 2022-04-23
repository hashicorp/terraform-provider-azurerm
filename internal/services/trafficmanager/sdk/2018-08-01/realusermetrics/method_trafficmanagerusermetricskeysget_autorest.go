package realusermetrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type TrafficManagerUserMetricsKeysGetResponse struct {
	HttpResponse *http.Response
	Model        *UserMetricsModel
}

// TrafficManagerUserMetricsKeysGet ...
func (c RealUserMetricsClient) TrafficManagerUserMetricsKeysGet(ctx context.Context, id commonids.SubscriptionId) (result TrafficManagerUserMetricsKeysGetResponse, err error) {
	req, err := c.preparerForTrafficManagerUserMetricsKeysGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "realusermetrics.RealUserMetricsClient", "TrafficManagerUserMetricsKeysGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "realusermetrics.RealUserMetricsClient", "TrafficManagerUserMetricsKeysGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTrafficManagerUserMetricsKeysGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "realusermetrics.RealUserMetricsClient", "TrafficManagerUserMetricsKeysGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTrafficManagerUserMetricsKeysGet prepares the TrafficManagerUserMetricsKeysGet request.
func (c RealUserMetricsClient) preparerForTrafficManagerUserMetricsKeysGet(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Network/trafficManagerUserMetricsKeys/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTrafficManagerUserMetricsKeysGet handles the response to the TrafficManagerUserMetricsKeysGet request. The method always
// closes the http.Response Body.
func (c RealUserMetricsClient) responderForTrafficManagerUserMetricsKeysGet(resp *http.Response) (result TrafficManagerUserMetricsKeysGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
