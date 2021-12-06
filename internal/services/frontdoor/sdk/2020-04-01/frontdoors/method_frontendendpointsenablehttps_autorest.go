package frontdoors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type FrontendEndpointsEnableHttpsResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// FrontendEndpointsEnableHttps ...
func (c FrontDoorsClient) FrontendEndpointsEnableHttps(ctx context.Context, id FrontendEndpointId, input CustomHttpsConfiguration) (result FrontendEndpointsEnableHttpsResponse, err error) {
	req, err := c.preparerForFrontendEndpointsEnableHttps(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsEnableHttps", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForFrontendEndpointsEnableHttps(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsEnableHttps", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// FrontendEndpointsEnableHttpsThenPoll performs FrontendEndpointsEnableHttps then polls until it's completed
func (c FrontDoorsClient) FrontendEndpointsEnableHttpsThenPoll(ctx context.Context, id FrontendEndpointId, input CustomHttpsConfiguration) error {
	result, err := c.FrontendEndpointsEnableHttps(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing FrontendEndpointsEnableHttps: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after FrontendEndpointsEnableHttps: %+v", err)
	}

	return nil
}

// preparerForFrontendEndpointsEnableHttps prepares the FrontendEndpointsEnableHttps request.
func (c FrontDoorsClient) preparerForFrontendEndpointsEnableHttps(ctx context.Context, id FrontendEndpointId, input CustomHttpsConfiguration) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/enableHttps", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForFrontendEndpointsEnableHttps sends the FrontendEndpointsEnableHttps request. The method will close the
// http.Response Body if it receives an error.
func (c FrontDoorsClient) senderForFrontendEndpointsEnableHttps(ctx context.Context, req *http.Request) (future FrontendEndpointsEnableHttpsResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
