package managementlocks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetAtSubscriptionLevelOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagementLockObject
}

// GetAtSubscriptionLevel ...
func (c ManagementLocksClient) GetAtSubscriptionLevel(ctx context.Context, id LockId) (result GetAtSubscriptionLevelOperationResponse, err error) {
	req, err := c.preparerForGetAtSubscriptionLevel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "GetAtSubscriptionLevel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "GetAtSubscriptionLevel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetAtSubscriptionLevel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "GetAtSubscriptionLevel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetAtSubscriptionLevel prepares the GetAtSubscriptionLevel request.
func (c ManagementLocksClient) preparerForGetAtSubscriptionLevel(ctx context.Context, id LockId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetAtSubscriptionLevel handles the response to the GetAtSubscriptionLevel request. The method always
// closes the http.Response Body.
func (c ManagementLocksClient) responderForGetAtSubscriptionLevel(resp *http.Response) (result GetAtSubscriptionLevelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
