package actiongroupsapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *ActionGroupList
}

// ActionGroupsListByResourceGroup ...
func (c ActionGroupsAPIsClient) ActionGroupsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result ActionGroupsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForActionGroupsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListByResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActionGroupsListByResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListByResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActionGroupsListByResourceGroup prepares the ActionGroupsListByResourceGroup request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/actionGroups", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForActionGroupsListByResourceGroup handles the response to the ActionGroupsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ActionGroupsAPIsClient) responderForActionGroupsListByResourceGroup(resp *http.Response) (result ActionGroupsListByResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
