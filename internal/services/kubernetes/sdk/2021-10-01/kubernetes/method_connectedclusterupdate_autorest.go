package kubernetes

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConnectedClusterUpdateResponse struct {
	HttpResponse *http.Response
	Model        *ConnectedCluster
}

// ConnectedClusterUpdate ...
func (c KubernetesClient) ConnectedClusterUpdate(ctx context.Context, id ConnectedClusterId, input ConnectedClusterPatch) (result ConnectedClusterUpdateResponse, err error) {
	req, err := c.preparerForConnectedClusterUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetes.KubernetesClient", "ConnectedClusterUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetes.KubernetesClient", "ConnectedClusterUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedClusterUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetes.KubernetesClient", "ConnectedClusterUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedClusterUpdate prepares the ConnectedClusterUpdate request.
func (c KubernetesClient) preparerForConnectedClusterUpdate(ctx context.Context, id ConnectedClusterId, input ConnectedClusterPatch) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConnectedClusterUpdate handles the response to the ConnectedClusterUpdate request. The method always
// closes the http.Response Body.
func (c KubernetesClient) responderForConnectedClusterUpdate(resp *http.Response) (result ConnectedClusterUpdateResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
