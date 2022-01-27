package checknameavailabilitywithsubscription

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CheckNameAvailabilityWithSubscriptionResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityOutput
}

// CheckNameAvailabilityWithSubscription ...
func (c CheckNameAvailabilityWithSubscriptionClient) CheckNameAvailabilityWithSubscription(ctx context.Context, id SubscriptionId, input CheckNameAvailabilityInput) (result CheckNameAvailabilityWithSubscriptionResponse, err error) {
	req, err := c.preparerForCheckNameAvailabilityWithSubscription(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailabilitywithsubscription.CheckNameAvailabilityWithSubscriptionClient", "CheckNameAvailabilityWithSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailabilitywithsubscription.CheckNameAvailabilityWithSubscriptionClient", "CheckNameAvailabilityWithSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckNameAvailabilityWithSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailabilitywithsubscription.CheckNameAvailabilityWithSubscriptionClient", "CheckNameAvailabilityWithSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckNameAvailabilityWithSubscription prepares the CheckNameAvailabilityWithSubscription request.
func (c CheckNameAvailabilityWithSubscriptionClient) preparerForCheckNameAvailabilityWithSubscription(ctx context.Context, id SubscriptionId, input CheckNameAvailabilityInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CDN/checkNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckNameAvailabilityWithSubscription handles the response to the CheckNameAvailabilityWithSubscription request. The method always
// closes the http.Response Body.
func (c CheckNameAvailabilityWithSubscriptionClient) responderForCheckNameAvailabilityWithSubscription(resp *http.Response) (result CheckNameAvailabilityWithSubscriptionResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
