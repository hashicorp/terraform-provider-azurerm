package hostpool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetrieveRegistrationTokenOperationResponse struct {
	HttpResponse *http.Response
	Model        *RegistrationInfo
}

// RetrieveRegistrationToken ...
func (c HostPoolClient) RetrieveRegistrationToken(ctx context.Context, id HostPoolId) (result RetrieveRegistrationTokenOperationResponse, err error) {
	req, err := c.preparerForRetrieveRegistrationToken(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hostpool.HostPoolClient", "RetrieveRegistrationToken", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "hostpool.HostPoolClient", "RetrieveRegistrationToken", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRetrieveRegistrationToken(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hostpool.HostPoolClient", "RetrieveRegistrationToken", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRetrieveRegistrationToken prepares the RetrieveRegistrationToken request.
func (c HostPoolClient) preparerForRetrieveRegistrationToken(ctx context.Context, id HostPoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/retrieveRegistrationToken", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRetrieveRegistrationToken handles the response to the RetrieveRegistrationToken request. The method always
// closes the http.Response Body.
func (c HostPoolClient) responderForRetrieveRegistrationToken(resp *http.Response) (result RetrieveRegistrationTokenOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
