package clusterextensions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

type ExtensionsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ExtensionsUpdate ...
func (c ClusterExtensionsClient) ExtensionsUpdate(ctx context.Context, id ExtensionId, input PatchExtension) (result ExtensionsUpdateOperationResponse, err error) {
	req, err := c.preparerForExtensionsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExtensionsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusterextensions.ClusterExtensionsClient", "ExtensionsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExtensionsUpdateThenPoll performs ExtensionsUpdate then polls until it's completed
func (c ClusterExtensionsClient) ExtensionsUpdateThenPoll(ctx context.Context, id ExtensionId, input PatchExtension) error {
	result, err := c.ExtensionsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ExtensionsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ExtensionsUpdate: %+v", err)
	}

	return nil
}

// preparerForExtensionsUpdate prepares the ExtensionsUpdate request.
func (c ClusterExtensionsClient) preparerForExtensionsUpdate(ctx context.Context, id ExtensionId, input PatchExtension) (*http.Request, error) {
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

// senderForExtensionsUpdate sends the ExtensionsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ClusterExtensionsClient) senderForExtensionsUpdate(ctx context.Context, req *http.Request) (future ExtensionsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
