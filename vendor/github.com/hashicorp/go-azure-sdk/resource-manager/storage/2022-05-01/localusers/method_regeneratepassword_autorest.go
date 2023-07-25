package localusers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegeneratePasswordOperationResponse struct {
	HttpResponse *http.Response
	Model        *LocalUserRegeneratePasswordResult
}

// RegeneratePassword ...
func (c LocalUsersClient) RegeneratePassword(ctx context.Context, id LocalUserId) (result RegeneratePasswordOperationResponse, err error) {
	req, err := c.preparerForRegeneratePassword(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "localusers.LocalUsersClient", "RegeneratePassword", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "localusers.LocalUsersClient", "RegeneratePassword", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRegeneratePassword(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "localusers.LocalUsersClient", "RegeneratePassword", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRegeneratePassword prepares the RegeneratePassword request.
func (c LocalUsersClient) preparerForRegeneratePassword(ctx context.Context, id LocalUserId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regeneratePassword", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRegeneratePassword handles the response to the RegeneratePassword request. The method always
// closes the http.Response Body.
func (c LocalUsersClient) responderForRegeneratePassword(resp *http.Response) (result RegeneratePasswordOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
