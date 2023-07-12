package actiongroupsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ActionGroupResource
}

// ActionGroupsGet ...
func (c ActionGroupsAPIsClient) ActionGroupsGet(ctx context.Context, id ActionGroupId) (result ActionGroupsGetOperationResponse, err error) {
	req, err := c.preparerForActionGroupsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActionGroupsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActionGroupsGet prepares the ActionGroupsGet request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsGet(ctx context.Context, id ActionGroupId) (*http.Request, error) {
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

// responderForActionGroupsGet handles the response to the ActionGroupsGet request. The method always
// closes the http.Response Body.
func (c ActionGroupsAPIsClient) responderForActionGroupsGet(resp *http.Response) (result ActionGroupsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
