package checkfrontdoornameavailability

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CheckFrontDoorNameAvailabilityResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityOutput
}

// CheckFrontDoorNameAvailability ...
func (c CheckFrontDoorNameAvailabilityClient) CheckFrontDoorNameAvailability(ctx context.Context, input CheckNameAvailabilityInput) (result CheckFrontDoorNameAvailabilityResponse, err error) {
	req, err := c.preparerForCheckFrontDoorNameAvailability(ctx, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checkfrontdoornameavailability.CheckFrontDoorNameAvailabilityClient", "CheckFrontDoorNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "checkfrontdoornameavailability.CheckFrontDoorNameAvailabilityClient", "CheckFrontDoorNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckFrontDoorNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checkfrontdoornameavailability.CheckFrontDoorNameAvailabilityClient", "CheckFrontDoorNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckFrontDoorNameAvailability prepares the CheckFrontDoorNameAvailability request.
func (c CheckFrontDoorNameAvailabilityClient) preparerForCheckFrontDoorNameAvailability(ctx context.Context, input CheckNameAvailabilityInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Network/checkFrontDoorNameAvailability"),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckFrontDoorNameAvailability handles the response to the CheckFrontDoorNameAvailability request. The method always
// closes the http.Response Body.
func (c CheckFrontDoorNameAvailabilityClient) responderForCheckFrontDoorNameAvailability(resp *http.Response) (result CheckFrontDoorNameAvailabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
