package locations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type GetCapabilityResponse struct {
	HttpResponse *http.Response
	Model        *CapabilityInformation
}

// GetCapability ...
func (c LocationsClient) GetCapability(ctx context.Context, id LocationId) (result GetCapabilityResponse, err error) {
	req, err := c.preparerForGetCapability(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "locations.LocationsClient", "GetCapability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "locations.LocationsClient", "GetCapability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetCapability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "locations.LocationsClient", "GetCapability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetCapability prepares the GetCapability request.
func (c LocationsClient) preparerForGetCapability(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/capability", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetCapability handles the response to the GetCapability request. The method always
// closes the http.Response Body.
func (c LocationsClient) responderForGetCapability(resp *http.Response) (result GetCapabilityResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
