package locations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CheckTrialAvailabilityResponse struct {
	HttpResponse *http.Response
	Model        *Trial
}

// CheckTrialAvailability ...
func (c LocationsClient) CheckTrialAvailability(ctx context.Context, id LocationId) (result CheckTrialAvailabilityResponse, err error) {
	req, err := c.preparerForCheckTrialAvailability(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "locations.LocationsClient", "CheckTrialAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "locations.LocationsClient", "CheckTrialAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckTrialAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "locations.LocationsClient", "CheckTrialAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckTrialAvailability prepares the CheckTrialAvailability request.
func (c LocationsClient) preparerForCheckTrialAvailability(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkTrialAvailability", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckTrialAvailability handles the response to the CheckTrialAvailability request. The method always
// closes the http.Response Body.
func (c LocationsClient) responderForCheckTrialAvailability(resp *http.Response) (result CheckTrialAvailabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
