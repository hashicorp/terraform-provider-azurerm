package scheduledactions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecuteByScopeOperationResponse struct {
	HttpResponse *http.Response
}

// ExecuteByScope ...
func (c ScheduledActionsClient) ExecuteByScope(ctx context.Context, id ScopedScheduledActionId) (result ExecuteByScopeOperationResponse, err error) {
	req, err := c.preparerForExecuteByScope(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "ExecuteByScope", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "ExecuteByScope", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForExecuteByScope(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "ExecuteByScope", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForExecuteByScope prepares the ExecuteByScope request.
func (c ScheduledActionsClient) preparerForExecuteByScope(ctx context.Context, id ScopedScheduledActionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/execute", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForExecuteByScope handles the response to the ExecuteByScope request. The method always
// closes the http.Response Body.
func (c ScheduledActionsClient) responderForExecuteByScope(resp *http.Response) (result ExecuteByScopeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
