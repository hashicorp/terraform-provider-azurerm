package eventhubsclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type ClustersCreateOrUpdateResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ClustersCreateOrUpdate ...
func (c EventHubsClustersClient) ClustersCreateOrUpdate(ctx context.Context, id ClusterId, input Cluster) (result ClustersCreateOrUpdateResponse, err error) {
	req, err := c.preparerForClustersCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForClustersCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ClustersCreateOrUpdateThenPoll performs ClustersCreateOrUpdate then polls until it's completed
func (c EventHubsClustersClient) ClustersCreateOrUpdateThenPoll(ctx context.Context, id ClusterId, input Cluster) error {
	result, err := c.ClustersCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ClustersCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ClustersCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForClustersCreateOrUpdate prepares the ClustersCreateOrUpdate request.
func (c EventHubsClustersClient) preparerForClustersCreateOrUpdate(ctx context.Context, id ClusterId, input Cluster) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForClustersCreateOrUpdate sends the ClustersCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c EventHubsClustersClient) senderForClustersCreateOrUpdate(ctx context.Context, req *http.Request) (future ClustersCreateOrUpdateResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
