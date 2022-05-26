package flux

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type ConfigurationsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ConfigurationsUpdate ...
func (c FluxClient) ConfigurationsUpdate(ctx context.Context, id FluxConfigurationId, input FluxConfigurationPatch) (result ConfigurationsUpdateOperationResponse, err error) {
	req, err := c.preparerForConfigurationsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "flux.FluxClient", "ConfigurationsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForConfigurationsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "flux.FluxClient", "ConfigurationsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ConfigurationsUpdateThenPoll performs ConfigurationsUpdate then polls until it's completed
func (c FluxClient) ConfigurationsUpdateThenPoll(ctx context.Context, id FluxConfigurationId, input FluxConfigurationPatch) error {
	result, err := c.ConfigurationsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ConfigurationsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ConfigurationsUpdate: %+v", err)
	}

	return nil
}

// preparerForConfigurationsUpdate prepares the ConfigurationsUpdate request.
func (c FluxClient) preparerForConfigurationsUpdate(ctx context.Context, id FluxConfigurationId, input FluxConfigurationPatch) (*http.Request, error) {
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

// senderForConfigurationsUpdate sends the ConfigurationsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c FluxClient) senderForConfigurationsUpdate(ctx context.Context, req *http.Request) (future ConfigurationsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
