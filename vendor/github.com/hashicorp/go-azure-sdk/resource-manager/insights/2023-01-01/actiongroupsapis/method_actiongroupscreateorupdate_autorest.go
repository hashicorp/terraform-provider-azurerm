package actiongroupsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ActionGroupResource
}

// ActionGroupsCreateOrUpdate ...
func (c ActionGroupsAPIsClient) ActionGroupsCreateOrUpdate(ctx context.Context, id ActionGroupId, input ActionGroupResource) (result ActionGroupsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForActionGroupsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActionGroupsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActionGroupsCreateOrUpdate prepares the ActionGroupsCreateOrUpdate request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsCreateOrUpdate(ctx context.Context, id ActionGroupId, input ActionGroupResource) (*http.Request, error) {
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

// responderForActionGroupsCreateOrUpdate handles the response to the ActionGroupsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c ActionGroupsAPIsClient) responderForActionGroupsCreateOrUpdate(resp *http.Response) (result ActionGroupsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
