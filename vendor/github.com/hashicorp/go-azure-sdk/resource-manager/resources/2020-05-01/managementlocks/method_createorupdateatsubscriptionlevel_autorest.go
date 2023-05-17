package managementlocks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateOrUpdateAtSubscriptionLevelOperationResponse struct {
	HttpResponse *http.Response
	Model        *ManagementLockObject
}

// CreateOrUpdateAtSubscriptionLevel ...
func (c ManagementLocksClient) CreateOrUpdateAtSubscriptionLevel(ctx context.Context, id LockId, input ManagementLockObject) (result CreateOrUpdateAtSubscriptionLevelOperationResponse, err error) {
	req, err := c.preparerForCreateOrUpdateAtSubscriptionLevel(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "CreateOrUpdateAtSubscriptionLevel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "CreateOrUpdateAtSubscriptionLevel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateOrUpdateAtSubscriptionLevel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "CreateOrUpdateAtSubscriptionLevel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateOrUpdateAtSubscriptionLevel prepares the CreateOrUpdateAtSubscriptionLevel request.
func (c ManagementLocksClient) preparerForCreateOrUpdateAtSubscriptionLevel(ctx context.Context, id LockId, input ManagementLockObject) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCreateOrUpdateAtSubscriptionLevel handles the response to the CreateOrUpdateAtSubscriptionLevel request. The method always
// closes the http.Response Body.
func (c ManagementLocksClient) responderForCreateOrUpdateAtSubscriptionLevel(resp *http.Response) (result CreateOrUpdateAtSubscriptionLevelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
