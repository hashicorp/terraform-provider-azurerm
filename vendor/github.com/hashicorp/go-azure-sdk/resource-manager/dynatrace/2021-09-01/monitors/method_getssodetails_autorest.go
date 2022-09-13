package monitors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetSSODetailsOperationResponse struct {
	HttpResponse *http.Response
	Model        *SSODetailsResponse
}

// GetSSODetails ...
func (c MonitorsClient) GetSSODetails(ctx context.Context, id MonitorId, input SSODetailsRequest) (result GetSSODetailsOperationResponse, err error) {
	req, err := c.preparerForGetSSODetails(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "GetSSODetails", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "GetSSODetails", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetSSODetails(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "GetSSODetails", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetSSODetails prepares the GetSSODetails request.
func (c MonitorsClient) preparerForGetSSODetails(ctx context.Context, id MonitorId, input SSODetailsRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getSSODetails", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetSSODetails handles the response to the GetSSODetails request. The method always
// closes the http.Response Body.
func (c MonitorsClient) responderForGetSSODetails(resp *http.Response) (result GetSSODetailsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
