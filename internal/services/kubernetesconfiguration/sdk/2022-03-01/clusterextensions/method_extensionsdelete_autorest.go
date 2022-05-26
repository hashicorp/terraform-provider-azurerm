package clusterextensions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type ExtensionsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type ExtensionsDeleteOperationOptions struct {
	ForceDelete *bool
}

func DefaultExtensionsDeleteOperationOptions() ExtensionsDeleteOperationOptions {
	return ExtensionsDeleteOperationOptions{}
}

func (o ExtensionsDeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ExtensionsDeleteOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ForceDelete != nil {
		out["forceDelete"] = *o.ForceDelete
	}

	return out
}

// ExtensionsDelete ...
func (c ClusterExtensionsClient) ExtensionsDelete(ctx context.Context, id ExtensionId, options ExtensionsDeleteOperationOptions) (result ExtensionsDeleteOperationResponse, err error) {
	req, err := c.preparerForExtensionsDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExtensionsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExtensionsDeleteThenPoll performs ExtensionsDelete then polls until it's completed
func (c ClusterExtensionsClient) ExtensionsDeleteThenPoll(ctx context.Context, id ExtensionId, options ExtensionsDeleteOperationOptions) error {
	result, err := c.ExtensionsDelete(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing ExtensionsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ExtensionsDelete: %+v", err)
	}

	return nil
}

// preparerForExtensionsDelete prepares the ExtensionsDelete request.
func (c ClusterExtensionsClient) preparerForExtensionsDelete(ctx context.Context, id ExtensionId, options ExtensionsDeleteOperationOptions) (*http.Request, error) {
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

// senderForExtensionsDelete sends the ExtensionsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c ClusterExtensionsClient) senderForExtensionsDelete(ctx context.Context, req *http.Request) (future ExtensionsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
