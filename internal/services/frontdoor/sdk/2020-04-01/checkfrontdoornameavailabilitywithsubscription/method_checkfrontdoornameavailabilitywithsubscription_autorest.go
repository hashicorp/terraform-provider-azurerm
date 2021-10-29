package checkfrontdoornameavailabilitywithsubscription

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CheckFrontDoorNameAvailabilityWithSubscriptionResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityOutput
}

// CheckFrontDoorNameAvailabilityWithSubscription ...
func (c CheckFrontDoorNameAvailabilityWithSubscriptionClient) CheckFrontDoorNameAvailabilityWithSubscription(ctx context.Context, id SubscriptionId, input CheckNameAvailabilityInput) (result CheckFrontDoorNameAvailabilityWithSubscriptionResponse, err error) {
	req, err := c.preparerForCheckFrontDoorNameAvailabilityWithSubscription(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checkfrontdoornameavailabilitywithsubscription.CheckFrontDoorNameAvailabilityWithSubscriptionClient", "CheckFrontDoorNameAvailabilityWithSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "checkfrontdoornameavailabilitywithsubscription.CheckFrontDoorNameAvailabilityWithSubscriptionClient", "CheckFrontDoorNameAvailabilityWithSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckFrontDoorNameAvailabilityWithSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checkfrontdoornameavailabilitywithsubscription.CheckFrontDoorNameAvailabilityWithSubscriptionClient", "CheckFrontDoorNameAvailabilityWithSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckFrontDoorNameAvailabilityWithSubscription prepares the CheckFrontDoorNameAvailabilityWithSubscription request.
func (c CheckFrontDoorNameAvailabilityWithSubscriptionClient) preparerForCheckFrontDoorNameAvailabilityWithSubscription(ctx context.Context, id SubscriptionId, input CheckNameAvailabilityInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Network/checkFrontDoorNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckFrontDoorNameAvailabilityWithSubscription handles the response to the CheckFrontDoorNameAvailabilityWithSubscription request. The method always
// closes the http.Response Body.
func (c CheckFrontDoorNameAvailabilityWithSubscriptionClient) responderForCheckFrontDoorNameAvailabilityWithSubscription(resp *http.Response) (result CheckFrontDoorNameAvailabilityWithSubscriptionResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
