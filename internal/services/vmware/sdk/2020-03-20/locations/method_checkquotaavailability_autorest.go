package locations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type CheckQuotaAvailabilityResponse struct {
	HttpResponse *http.Response
	Model        *Quota
}

// CheckQuotaAvailability ...
func (c LocationsClient) CheckQuotaAvailability(ctx context.Context, id LocationId) (result CheckQuotaAvailabilityResponse, err error) {
	req, err := c.preparerForCheckQuotaAvailability(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "locations.LocationsClient", "CheckQuotaAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "locations.LocationsClient", "CheckQuotaAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckQuotaAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "locations.LocationsClient", "CheckQuotaAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckQuotaAvailability prepares the CheckQuotaAvailability request.
func (c LocationsClient) preparerForCheckQuotaAvailability(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkQuotaAvailability", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckQuotaAvailability handles the response to the CheckQuotaAvailability request. The method always
// closes the http.Response Body.
func (c LocationsClient) responderForCheckQuotaAvailability(resp *http.Response) (result CheckQuotaAvailabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
