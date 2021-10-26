package frontdoors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type FrontendEndpointsDisableHttpsResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// FrontendEndpointsDisableHttps ...
func (c FrontDoorsClient) FrontendEndpointsDisableHttps(ctx context.Context, id FrontendEndpointId) (result FrontendEndpointsDisableHttpsResponse, err error) {
	req, err := c.preparerForFrontendEndpointsDisableHttps(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsDisableHttps", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForFrontendEndpointsDisableHttps(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsDisableHttps", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// FrontendEndpointsDisableHttpsThenPoll performs FrontendEndpointsDisableHttps then polls until it's completed
func (c FrontDoorsClient) FrontendEndpointsDisableHttpsThenPoll(ctx context.Context, id FrontendEndpointId) error {
	result, err := c.FrontendEndpointsDisableHttps(ctx, id)
	if err != nil {
		return fmt.Errorf("performing FrontendEndpointsDisableHttps: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after FrontendEndpointsDisableHttps: %+v", err)
	}

	return nil
}

// preparerForFrontendEndpointsDisableHttps prepares the FrontendEndpointsDisableHttps request.
func (c FrontDoorsClient) preparerForFrontendEndpointsDisableHttps(ctx context.Context, id FrontendEndpointId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/disableHttps", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForFrontendEndpointsDisableHttps sends the FrontendEndpointsDisableHttps request. The method will close the
// http.Response Body if it receives an error.
func (c FrontDoorsClient) senderForFrontendEndpointsDisableHttps(ctx context.Context, req *http.Request) (future FrontendEndpointsDisableHttpsResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
