package autoscalevcores

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListBySubscriptionResponse struct {
	HttpResponse *http.Response
	Model        *AutoScaleVCoreListResult
}

// ListBySubscription ...
func (c AutoScaleVCoresClient) ListBySubscription(ctx context.Context, id SubscriptionId) (result ListBySubscriptionResponse, err error) {
	req, err := c.preparerForListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "autoscalevcores.AutoScaleVCoresClient", "ListBySubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "autoscalevcores.AutoScaleVCoresClient", "ListBySubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListBySubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "autoscalevcores.AutoScaleVCoresClient", "ListBySubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListBySubscription prepares the ListBySubscription request.
func (c AutoScaleVCoresClient) preparerForListBySubscription(ctx context.Context, id SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.PowerBIDedicated/autoScaleVCores", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListBySubscription handles the response to the ListBySubscription request. The method always
// closes the http.Response Body.
func (c AutoScaleVCoresClient) responderForListBySubscription(resp *http.Response) (result ListBySubscriptionResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
