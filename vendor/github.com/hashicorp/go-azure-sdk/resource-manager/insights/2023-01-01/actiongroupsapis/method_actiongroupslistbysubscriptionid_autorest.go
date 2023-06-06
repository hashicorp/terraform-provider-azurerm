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

type ActionGroupsListBySubscriptionIdOperationResponse struct {
	HttpResponse *http.Response
	Model        *ActionGroupList
}

// ActionGroupsListBySubscriptionId ...
func (c ActionGroupsAPIsClient) ActionGroupsListBySubscriptionId(ctx context.Context, id commonids.SubscriptionId) (result ActionGroupsListBySubscriptionIdOperationResponse, err error) {
	req, err := c.preparerForActionGroupsListBySubscriptionId(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListBySubscriptionId", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListBySubscriptionId", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActionGroupsListBySubscriptionId(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListBySubscriptionId", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActionGroupsListBySubscriptionId prepares the ActionGroupsListBySubscriptionId request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsListBySubscriptionId(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
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

// responderForActionGroupsListBySubscriptionId handles the response to the ActionGroupsListBySubscriptionId request. The method always
// closes the http.Response Body.
func (c ActionGroupsAPIsClient) responderForActionGroupsListBySubscriptionId(resp *http.Response) (result ActionGroupsListBySubscriptionIdOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
