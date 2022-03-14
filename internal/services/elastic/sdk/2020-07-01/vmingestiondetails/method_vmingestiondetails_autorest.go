package vmingestiondetails

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type VMIngestionDetailsOperationResponse struct {
	HttpResponse *http.Response
	Model        *VMIngestionDetailsResponse
}

// VMIngestionDetails ...
func (c VMIngestionDetailsClient) VMIngestionDetails(ctx context.Context, id MonitorId) (result VMIngestionDetailsOperationResponse, err error) {
	req, err := c.preparerForVMIngestionDetails(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmingestiondetails.VMIngestionDetailsClient", "VMIngestionDetails", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmingestiondetails.VMIngestionDetailsClient", "VMIngestionDetails", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForVMIngestionDetails(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmingestiondetails.VMIngestionDetailsClient", "VMIngestionDetails", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForVMIngestionDetails prepares the VMIngestionDetails request.
func (c VMIngestionDetailsClient) preparerForVMIngestionDetails(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/vmIngestionDetails", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVMIngestionDetails handles the response to the VMIngestionDetails request. The method always
// closes the http.Response Body.
func (c VMIngestionDetailsClient) responderForVMIngestionDetails(resp *http.Response) (result VMIngestionDetailsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
