package eventhubsclustersavailableclusterregions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ClustersListAvailableClusterRegionResponse struct {
	HttpResponse *http.Response
	Model        *AvailableClustersList
}

// ClustersListAvailableClusterRegion ...
func (c EventHubsClustersAvailableClusterRegionsClient) ClustersListAvailableClusterRegion(ctx context.Context, id SubscriptionId) (result ClustersListAvailableClusterRegionResponse, err error) {
	req, err := c.preparerForClustersListAvailableClusterRegion(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersavailableclusterregions.EventHubsClustersAvailableClusterRegionsClient", "ClustersListAvailableClusterRegion", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersavailableclusterregions.EventHubsClustersAvailableClusterRegionsClient", "ClustersListAvailableClusterRegion", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForClustersListAvailableClusterRegion(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersavailableclusterregions.EventHubsClustersAvailableClusterRegionsClient", "ClustersListAvailableClusterRegion", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForClustersListAvailableClusterRegion prepares the ClustersListAvailableClusterRegion request.
func (c EventHubsClustersAvailableClusterRegionsClient) preparerForClustersListAvailableClusterRegion(ctx context.Context, id SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.EventHub/availableClusterRegions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForClustersListAvailableClusterRegion handles the response to the ClustersListAvailableClusterRegion request. The method always
// closes the http.Response Body.
func (c EventHubsClustersAvailableClusterRegionsClient) responderForClustersListAvailableClusterRegion(resp *http.Response) (result ClustersListAvailableClusterRegionResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
