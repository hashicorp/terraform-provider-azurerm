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

type GetAccountCredentialsOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccountInfoSecure
}

// GetAccountCredentials ...
func (c MonitorsClient) GetAccountCredentials(ctx context.Context, id MonitorId) (result GetAccountCredentialsOperationResponse, err error) {
	req, err := c.preparerForGetAccountCredentials(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "GetAccountCredentials", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "GetAccountCredentials", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetAccountCredentials(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "GetAccountCredentials", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetAccountCredentials prepares the GetAccountCredentials request.
func (c MonitorsClient) preparerForGetAccountCredentials(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getAccountCredentials", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetAccountCredentials handles the response to the GetAccountCredentials request. The method always
// closes the http.Response Body.
func (c MonitorsClient) responderForGetAccountCredentials(resp *http.Response) (result GetAccountCredentialsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
