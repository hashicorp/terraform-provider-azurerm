package checkfrontdoornameavailabilitywithsubscription

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type FrontDoorNameAvailabilityWithSubscriptionCheckResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityOutput
}

// FrontDoorNameAvailabilityWithSubscriptionCheck ...
func (c CheckFrontDoorNameAvailabilityWithSubscriptionClient) FrontDoorNameAvailabilityWithSubscriptionCheck(ctx context.Context, id SubscriptionId, input CheckNameAvailabilityInput) (result FrontDoorNameAvailabilityWithSubscriptionCheckResponse, err error) {
	req, err := c.preparerForFrontDoorNameAvailabilityWithSubscriptionCheck(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checkfrontdoornameavailabilitywithsubscription.CheckFrontDoorNameAvailabilityWithSubscriptionClient", "FrontDoorNameAvailabilityWithSubscriptionCheck", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "checkfrontdoornameavailabilitywithsubscription.CheckFrontDoorNameAvailabilityWithSubscriptionClient", "FrontDoorNameAvailabilityWithSubscriptionCheck", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForFrontDoorNameAvailabilityWithSubscriptionCheck(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checkfrontdoornameavailabilitywithsubscription.CheckFrontDoorNameAvailabilityWithSubscriptionClient", "FrontDoorNameAvailabilityWithSubscriptionCheck", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForFrontDoorNameAvailabilityWithSubscriptionCheck prepares the FrontDoorNameAvailabilityWithSubscriptionCheck request.
func (c CheckFrontDoorNameAvailabilityWithSubscriptionClient) preparerForFrontDoorNameAvailabilityWithSubscriptionCheck(ctx context.Context, id SubscriptionId, input CheckNameAvailabilityInput) (*http.Request, error) {
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

// responderForFrontDoorNameAvailabilityWithSubscriptionCheck handles the response to the FrontDoorNameAvailabilityWithSubscriptionCheck request. The method always
// closes the http.Response Body.
func (c CheckFrontDoorNameAvailabilityWithSubscriptionClient) responderForFrontDoorNameAvailabilityWithSubscriptionCheck(resp *http.Response) (result FrontDoorNameAvailabilityWithSubscriptionCheckResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
