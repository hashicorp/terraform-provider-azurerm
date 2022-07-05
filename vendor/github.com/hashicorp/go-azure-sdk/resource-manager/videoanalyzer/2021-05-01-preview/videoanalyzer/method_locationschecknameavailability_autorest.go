package videoanalyzer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationsCheckNameAvailabilityOperationResponse struct {
	HttpResponse *http.Response
	Model        *CheckNameAvailabilityResponse
}

// LocationsCheckNameAvailability ...
func (c VideoAnalyzerClient) LocationsCheckNameAvailability(ctx context.Context, id LocationId, input CheckNameAvailabilityRequest) (result LocationsCheckNameAvailabilityOperationResponse, err error) {
	req, err := c.preparerForLocationsCheckNameAvailability(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "LocationsCheckNameAvailability", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "LocationsCheckNameAvailability", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLocationsCheckNameAvailability(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "LocationsCheckNameAvailability", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLocationsCheckNameAvailability prepares the LocationsCheckNameAvailability request.
func (c VideoAnalyzerClient) preparerForLocationsCheckNameAvailability(ctx context.Context, id LocationId, input CheckNameAvailabilityRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkNameAvailability", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLocationsCheckNameAvailability handles the response to the LocationsCheckNameAvailability request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForLocationsCheckNameAvailability(resp *http.Response) (result LocationsCheckNameAvailabilityOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
