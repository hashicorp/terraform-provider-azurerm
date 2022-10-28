package deviceupdates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountsHeadOperationResponse struct {
	HttpResponse *http.Response
}

// AccountsHead ...
func (c DeviceupdatesClient) AccountsHead(ctx context.Context, id AccountId) (result AccountsHeadOperationResponse, err error) {
	req, err := c.preparerForAccountsHead(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "AccountsHead", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "AccountsHead", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountsHead(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "AccountsHead", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountsHead prepares the AccountsHead request.
func (c DeviceupdatesClient) preparerForAccountsHead(ctx context.Context, id AccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsHead(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAccountsHead handles the response to the AccountsHead request. The method always
// closes the http.Response Body.
func (c DeviceupdatesClient) responderForAccountsHead(resp *http.Response) (result AccountsHeadOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
