package managementlocks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteAtSubscriptionLevelOperationResponse struct {
	HttpResponse *http.Response
}

// DeleteAtSubscriptionLevel ...
func (c ManagementLocksClient) DeleteAtSubscriptionLevel(ctx context.Context, id LockId) (result DeleteAtSubscriptionLevelOperationResponse, err error) {
	req, err := c.preparerForDeleteAtSubscriptionLevel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "DeleteAtSubscriptionLevel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "DeleteAtSubscriptionLevel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteAtSubscriptionLevel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "DeleteAtSubscriptionLevel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteAtSubscriptionLevel prepares the DeleteAtSubscriptionLevel request.
func (c ManagementLocksClient) preparerForDeleteAtSubscriptionLevel(ctx context.Context, id LockId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDeleteAtSubscriptionLevel handles the response to the DeleteAtSubscriptionLevel request. The method always
// closes the http.Response Body.
func (c ManagementLocksClient) responderForDeleteAtSubscriptionLevel(resp *http.Response) (result DeleteAtSubscriptionLevelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
