package eventhubsclusters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ClustersGetResponse struct {
	HttpResponse *http.Response
	Model        *Cluster
}

// ClustersGet ...
func (c EventHubsClustersClient) ClustersGet(ctx context.Context, id ClusterId) (result ClustersGetResponse, err error) {
	req, err := c.preparerForClustersGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForClustersGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForClustersGet prepares the ClustersGet request.
func (c EventHubsClustersClient) preparerForClustersGet(ctx context.Context, id ClusterId) (*http.Request, error) {
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

// responderForClustersGet handles the response to the ClustersGet request. The method always
// closes the http.Response Body.
func (c EventHubsClustersClient) responderForClustersGet(resp *http.Response) (result ClustersGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
