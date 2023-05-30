package actiongroupsapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsGetTestNotificationsAtActionGroupResourceLevelOperationResponse struct {
	HttpResponse *http.Response
	Model        *TestNotificationDetailsResponse
}

// ActionGroupsGetTestNotificationsAtActionGroupResourceLevel ...
func (c ActionGroupsAPIsClient) ActionGroupsGetTestNotificationsAtActionGroupResourceLevel(ctx context.Context, id NotificationStatusId) (result ActionGroupsGetTestNotificationsAtActionGroupResourceLevelOperationResponse, err error) {
	req, err := c.preparerForActionGroupsGetTestNotificationsAtActionGroupResourceLevel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsGetTestNotificationsAtActionGroupResourceLevel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsGetTestNotificationsAtActionGroupResourceLevel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActionGroupsGetTestNotificationsAtActionGroupResourceLevel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsGetTestNotificationsAtActionGroupResourceLevel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActionGroupsGetTestNotificationsAtActionGroupResourceLevel prepares the ActionGroupsGetTestNotificationsAtActionGroupResourceLevel request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsGetTestNotificationsAtActionGroupResourceLevel(ctx context.Context, id NotificationStatusId) (*http.Request, error) {
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

// responderForActionGroupsGetTestNotificationsAtActionGroupResourceLevel handles the response to the ActionGroupsGetTestNotificationsAtActionGroupResourceLevel request. The method always
// closes the http.Response Body.
func (c ActionGroupsAPIsClient) responderForActionGroupsGetTestNotificationsAtActionGroupResourceLevel(resp *http.Response) (result ActionGroupsGetTestNotificationsAtActionGroupResourceLevelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
