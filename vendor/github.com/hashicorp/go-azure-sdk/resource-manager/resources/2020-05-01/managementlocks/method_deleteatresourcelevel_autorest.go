package managementlocks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteAtResourceLevelOperationResponse struct {
	HttpResponse *http.Response
}

// DeleteAtResourceLevel ...
func (c ManagementLocksClient) DeleteAtResourceLevel(ctx context.Context, id ResourceLockId) (result DeleteAtResourceLevelOperationResponse, err error) {
	req, err := c.preparerForDeleteAtResourceLevel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "DeleteAtResourceLevel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "DeleteAtResourceLevel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteAtResourceLevel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "DeleteAtResourceLevel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteAtResourceLevel prepares the DeleteAtResourceLevel request.
func (c ManagementLocksClient) preparerForDeleteAtResourceLevel(ctx context.Context, id ResourceLockId) (*http.Request, error) {
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

// responderForDeleteAtResourceLevel handles the response to the DeleteAtResourceLevel request. The method always
// closes the http.Response Body.
func (c ManagementLocksClient) responderForDeleteAtResourceLevel(resp *http.Response) (result DeleteAtResourceLevelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
