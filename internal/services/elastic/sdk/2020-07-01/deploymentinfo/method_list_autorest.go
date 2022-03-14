package deploymentinfo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListOperationResponse struct {
	HttpResponse *http.Response
	Model        *DeploymentInfoResponse
}

// List ...
func (c DeploymentInfoClient) List(ctx context.Context, id MonitorId) (result ListOperationResponse, err error) {
	req, err := c.preparerForList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentinfo.DeploymentInfoClient", "List", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentinfo.DeploymentInfoClient", "List", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentinfo.DeploymentInfoClient", "List", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForList prepares the List request.
func (c DeploymentInfoClient) preparerForList(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listDeploymentInfo", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForList handles the response to the List request. The method always
// closes the http.Response Body.
func (c DeploymentInfoClient) responderForList(resp *http.Response) (result ListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
