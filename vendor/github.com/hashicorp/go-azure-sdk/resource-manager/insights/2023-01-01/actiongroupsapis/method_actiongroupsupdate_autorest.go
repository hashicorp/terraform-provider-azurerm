package actiongroupsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ActionGroupResource
}

// ActionGroupsUpdate ...
func (c ActionGroupsAPIsClient) ActionGroupsUpdate(ctx context.Context, id ActionGroupId, input ActionGroupPatchBody) (result ActionGroupsUpdateOperationResponse, err error) {
	req, err := c.preparerForActionGroupsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActionGroupsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActionGroupsUpdate prepares the ActionGroupsUpdate request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsUpdate(ctx context.Context, id ActionGroupId, input ActionGroupPatchBody) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForActionGroupsUpdate handles the response to the ActionGroupsUpdate request. The method always
// closes the http.Response Body.
func (c ActionGroupsAPIsClient) responderForActionGroupsUpdate(resp *http.Response) (result ActionGroupsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
