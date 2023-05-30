package actiongroupsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// ActionGroupsDelete ...
func (c ActionGroupsAPIsClient) ActionGroupsDelete(ctx context.Context, id ActionGroupId) (result ActionGroupsDeleteOperationResponse, err error) {
	req, err := c.preparerForActionGroupsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActionGroupsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActionGroupsDelete prepares the ActionGroupsDelete request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsDelete(ctx context.Context, id ActionGroupId) (*http.Request, error) {
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

// responderForActionGroupsDelete handles the response to the ActionGroupsDelete request. The method always
// closes the http.Response Body.
func (c ActionGroupsAPIsClient) responderForActionGroupsDelete(resp *http.Response) (result ActionGroupsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
