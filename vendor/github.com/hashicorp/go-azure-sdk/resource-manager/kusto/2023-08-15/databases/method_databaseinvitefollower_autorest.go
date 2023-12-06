package databases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseInviteFollowerOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabaseInviteFollowerResult
}

// DatabaseInviteFollower ...
func (c DatabasesClient) DatabaseInviteFollower(ctx context.Context, id DatabaseId, input DatabaseInviteFollowerRequest) (result DatabaseInviteFollowerOperationResponse, err error) {
	req, err := c.preparerForDatabaseInviteFollower(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "DatabaseInviteFollower", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "DatabaseInviteFollower", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDatabaseInviteFollower(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "DatabaseInviteFollower", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDatabaseInviteFollower prepares the DatabaseInviteFollower request.
func (c DatabasesClient) preparerForDatabaseInviteFollower(ctx context.Context, id DatabaseId, input DatabaseInviteFollowerRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/inviteFollower", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDatabaseInviteFollower handles the response to the DatabaseInviteFollower request. The method always
// closes the http.Response Body.
func (c DatabasesClient) responderForDatabaseInviteFollower(resp *http.Response) (result DatabaseInviteFollowerOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
