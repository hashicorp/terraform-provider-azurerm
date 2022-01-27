package afdprofiles

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CheckHostNameAvailabilityResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityOutput
}

// CheckHostNameAvailability ...
func (c AFDProfilesClient) CheckHostNameAvailability(ctx context.Context, id ProfileId, input CheckHostNameAvailabilityInput) (result CheckHostNameAvailabilityResponse, err error) {
	req, err := c.preparerForCheckHostNameAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "afdprofiles.AFDProfilesClient", "CheckHostNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "afdprofiles.AFDProfilesClient", "CheckHostNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckHostNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "afdprofiles.AFDProfilesClient", "CheckHostNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckHostNameAvailability prepares the CheckHostNameAvailability request.
func (c AFDProfilesClient) preparerForCheckHostNameAvailability(ctx context.Context, id ProfileId, input CheckHostNameAvailabilityInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkHostNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckHostNameAvailability handles the response to the CheckHostNameAvailability request. The method always
// closes the http.Response Body.
func (c AFDProfilesClient) responderForCheckHostNameAvailability(resp *http.Response) (result CheckHostNameAvailabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
