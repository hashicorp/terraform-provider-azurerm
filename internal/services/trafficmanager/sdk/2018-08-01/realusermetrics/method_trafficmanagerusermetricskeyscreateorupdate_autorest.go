package realusermetrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type TrafficManagerUserMetricsKeysCreateOrUpdateResponse struct {
	HttpResponse *http.Response
	Model        *UserMetricsModel
}

// TrafficManagerUserMetricsKeysCreateOrUpdate ...
func (c RealUserMetricsClient) TrafficManagerUserMetricsKeysCreateOrUpdate(ctx context.Context, id commonids.SubscriptionId) (result TrafficManagerUserMetricsKeysCreateOrUpdateResponse, err error) {
	req, err := c.preparerForTrafficManagerUserMetricsKeysCreateOrUpdate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "realusermetrics.RealUserMetricsClient", "TrafficManagerUserMetricsKeysCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "realusermetrics.RealUserMetricsClient", "TrafficManagerUserMetricsKeysCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTrafficManagerUserMetricsKeysCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "realusermetrics.RealUserMetricsClient", "TrafficManagerUserMetricsKeysCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTrafficManagerUserMetricsKeysCreateOrUpdate prepares the TrafficManagerUserMetricsKeysCreateOrUpdate request.
func (c RealUserMetricsClient) preparerForTrafficManagerUserMetricsKeysCreateOrUpdate(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Network/trafficManagerUserMetricsKeys/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTrafficManagerUserMetricsKeysCreateOrUpdate handles the response to the TrafficManagerUserMetricsKeysCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c RealUserMetricsClient) responderForTrafficManagerUserMetricsKeysCreateOrUpdate(resp *http.Response) (result TrafficManagerUserMetricsKeysCreateOrUpdateResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
