package eventhubsclustersnamespace

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ClustersListNamespacesResponse struct {
	HttpResponse *http.Response
	Model        *EHNamespaceIdListResult
}

// ClustersListNamespaces ...
func (c EventHubsClustersNamespaceClient) ClustersListNamespaces(ctx context.Context, id ClusterId) (result ClustersListNamespacesResponse, err error) {
	req, err := c.preparerForClustersListNamespaces(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersnamespace.EventHubsClustersNamespaceClient", "ClustersListNamespaces", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersnamespace.EventHubsClustersNamespaceClient", "ClustersListNamespaces", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForClustersListNamespaces(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclustersnamespace.EventHubsClustersNamespaceClient", "ClustersListNamespaces", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForClustersListNamespaces prepares the ClustersListNamespaces request.
func (c EventHubsClustersNamespaceClient) preparerForClustersListNamespaces(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/namespaces", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForClustersListNamespaces handles the response to the ClustersListNamespaces request. The method always
// closes the http.Response Body.
func (c EventHubsClustersNamespaceClient) responderForClustersListNamespaces(resp *http.Response) (result ClustersListNamespacesResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
