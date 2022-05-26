package flux

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type ConfigurationsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type ConfigurationsDeleteOperationOptions struct {
	ForceDelete *bool
}

func DefaultConfigurationsDeleteOperationOptions() ConfigurationsDeleteOperationOptions {
	return ConfigurationsDeleteOperationOptions{}
}

func (o ConfigurationsDeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ConfigurationsDeleteOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ForceDelete != nil {
		out["forceDelete"] = *o.ForceDelete
	}

	return out
}

// ConfigurationsDelete ...
func (c FluxClient) ConfigurationsDelete(ctx context.Context, id FluxConfigurationId, options ConfigurationsDeleteOperationOptions) (result ConfigurationsDeleteOperationResponse, err error) {
	req, err := c.preparerForConfigurationsDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "flux.FluxClient", "ConfigurationsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForConfigurationsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "flux.FluxClient", "ConfigurationsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ConfigurationsDeleteThenPoll performs ConfigurationsDelete then polls until it's completed
func (c FluxClient) ConfigurationsDeleteThenPoll(ctx context.Context, id FluxConfigurationId, options ConfigurationsDeleteOperationOptions) error {
	result, err := c.ConfigurationsDelete(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing ConfigurationsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ConfigurationsDelete: %+v", err)
	}

	return nil
}

// preparerForConfigurationsDelete prepares the ConfigurationsDelete request.
func (c FluxClient) preparerForConfigurationsDelete(ctx context.Context, id FluxConfigurationId, options ConfigurationsDeleteOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForConfigurationsDelete sends the ConfigurationsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c FluxClient) senderForConfigurationsDelete(ctx context.Context, req *http.Request) (future ConfigurationsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
