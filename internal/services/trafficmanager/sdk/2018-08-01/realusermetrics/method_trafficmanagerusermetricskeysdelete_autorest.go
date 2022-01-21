package realusermetrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type TrafficManagerUserMetricsKeysDeleteResponse struct {
	HttpResponse *http.Response
	Model        *DeleteOperationResult
}

// TrafficManagerUserMetricsKeysDelete ...
func (c RealUserMetricsClient) TrafficManagerUserMetricsKeysDelete(ctx context.Context, id commonids.SubscriptionId) (result TrafficManagerUserMetricsKeysDeleteResponse, err error) {
	req, err := c.preparerForTrafficManagerUserMetricsKeysDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "realusermetrics.RealUserMetricsClient", "TrafficManagerUserMetricsKeysDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "realusermetrics.RealUserMetricsClient", "TrafficManagerUserMetricsKeysDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTrafficManagerUserMetricsKeysDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "realusermetrics.RealUserMetricsClient", "TrafficManagerUserMetricsKeysDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTrafficManagerUserMetricsKeysDelete prepares the TrafficManagerUserMetricsKeysDelete request.
func (c RealUserMetricsClient) preparerForTrafficManagerUserMetricsKeysDelete(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Network/trafficManagerUserMetricsKeys/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTrafficManagerUserMetricsKeysDelete handles the response to the TrafficManagerUserMetricsKeysDelete request. The method always
// closes the http.Response Body.
func (c RealUserMetricsClient) responderForTrafficManagerUserMetricsKeysDelete(resp *http.Response) (result TrafficManagerUserMetricsKeysDeleteResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
