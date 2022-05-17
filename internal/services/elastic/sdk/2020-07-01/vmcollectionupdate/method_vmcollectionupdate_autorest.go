package vmcollectionupdate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type VMCollectionUpdateOperationResponse struct {
	HttpResponse *http.Response
}

// VMCollectionUpdate ...
func (c VMCollectionUpdateClient) VMCollectionUpdate(ctx context.Context, id MonitorId, input VMCollectionUpdate) (result VMCollectionUpdateOperationResponse, err error) {
	req, err := c.preparerForVMCollectionUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmcollectionupdate.VMCollectionUpdateClient", "VMCollectionUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmcollectionupdate.VMCollectionUpdateClient", "VMCollectionUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVMCollectionUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmcollectionupdate.VMCollectionUpdateClient", "VMCollectionUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVMCollectionUpdate prepares the VMCollectionUpdate request.
func (c VMCollectionUpdateClient) preparerForVMCollectionUpdate(ctx context.Context, id MonitorId, input VMCollectionUpdate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/vmCollectionUpdate", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVMCollectionUpdate handles the response to the VMCollectionUpdate request. The method always
// closes the http.Response Body.
func (c VMCollectionUpdateClient) responderForVMCollectionUpdate(resp *http.Response) (result VMCollectionUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
