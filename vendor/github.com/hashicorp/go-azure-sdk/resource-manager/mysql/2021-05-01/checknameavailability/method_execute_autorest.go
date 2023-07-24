package checknameavailability

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecuteOperationResponse struct {
	HttpResponse *http.Response
	Model        *NameAvailability
}

// Execute ...
func (c CheckNameAvailabilityClient) Execute(ctx context.Context, id LocationId, input NameAvailabilityRequest) (result ExecuteOperationResponse, err error) {
	req, err := c.preparerForExecute(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailability.CheckNameAvailabilityClient", "Execute", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailability.CheckNameAvailabilityClient", "Execute", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForExecute(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "checknameavailability.CheckNameAvailabilityClient", "Execute", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForExecute prepares the Execute request.
func (c CheckNameAvailabilityClient) preparerForExecute(ctx context.Context, id LocationId, input NameAvailabilityRequest) (*http.Request, error) {
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

// responderForExecute handles the response to the Execute request. The method always
// closes the http.Response Body.
func (c CheckNameAvailabilityClient) responderForExecute(resp *http.Response) (result ExecuteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
