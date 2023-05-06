package actiongroupsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsGetTestNotificationsAtResourceGroupLevelOperationResponse struct {
	HttpResponse *http.Response
	Model        *TestNotificationDetailsResponse
}

// ActionGroupsGetTestNotificationsAtResourceGroupLevel ...
func (c ActionGroupsAPIsClient) ActionGroupsGetTestNotificationsAtResourceGroupLevel(ctx context.Context, id ProviderNotificationStatusId) (result ActionGroupsGetTestNotificationsAtResourceGroupLevelOperationResponse, err error) {
	req, err := c.preparerForActionGroupsGetTestNotificationsAtResourceGroupLevel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsGetTestNotificationsAtResourceGroupLevel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsGetTestNotificationsAtResourceGroupLevel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActionGroupsGetTestNotificationsAtResourceGroupLevel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsGetTestNotificationsAtResourceGroupLevel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActionGroupsGetTestNotificationsAtResourceGroupLevel prepares the ActionGroupsGetTestNotificationsAtResourceGroupLevel request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsGetTestNotificationsAtResourceGroupLevel(ctx context.Context, id ProviderNotificationStatusId) (*http.Request, error) {
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

// responderForActionGroupsGetTestNotificationsAtResourceGroupLevel handles the response to the ActionGroupsGetTestNotificationsAtResourceGroupLevel request. The method always
// closes the http.Response Body.
func (c ActionGroupsAPIsClient) responderForActionGroupsGetTestNotificationsAtResourceGroupLevel(resp *http.Response) (result ActionGroupsGetTestNotificationsAtResourceGroupLevelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
