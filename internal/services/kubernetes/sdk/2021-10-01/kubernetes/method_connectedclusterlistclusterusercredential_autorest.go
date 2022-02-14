package kubernetes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ConnectedClusterListClusterUserCredentialResponse struct {
	HttpResponse *http.Response
	Model        *CredentialResults
}

// ConnectedClusterListClusterUserCredential ...
func (c KubernetesClient) ConnectedClusterListClusterUserCredential(ctx context.Context, id ConnectedClusterId, input ListClusterUserCredentialProperties) (result ConnectedClusterListClusterUserCredentialResponse, err error) {
	req, err := c.preparerForConnectedClusterListClusterUserCredential(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetes.KubernetesClient", "ConnectedClusterListClusterUserCredential", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetes.KubernetesClient", "ConnectedClusterListClusterUserCredential", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedClusterListClusterUserCredential(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "kubernetes.KubernetesClient", "ConnectedClusterListClusterUserCredential", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedClusterListClusterUserCredential prepares the ConnectedClusterListClusterUserCredential request.
func (c KubernetesClient) preparerForConnectedClusterListClusterUserCredential(ctx context.Context, id ConnectedClusterId, input ListClusterUserCredentialProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listClusterUserCredential", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConnectedClusterListClusterUserCredential handles the response to the ConnectedClusterListClusterUserCredential request. The method always
// closes the http.Response Body.
func (c KubernetesClient) responderForConnectedClusterListClusterUserCredential(resp *http.Response) (result ConnectedClusterListClusterUserCredentialResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
