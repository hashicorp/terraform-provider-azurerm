package actiongroupsapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsCreateNotificationsAtActionGroupResourceLevelOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ActionGroupsCreateNotificationsAtActionGroupResourceLevel ...
func (c ActionGroupsAPIsClient) ActionGroupsCreateNotificationsAtActionGroupResourceLevel(ctx context.Context, id ActionGroupId, input NotificationRequestBody) (result ActionGroupsCreateNotificationsAtActionGroupResourceLevelOperationResponse, err error) {
	req, err := c.preparerForActionGroupsCreateNotificationsAtActionGroupResourceLevel(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsCreateNotificationsAtActionGroupResourceLevel", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForActionGroupsCreateNotificationsAtActionGroupResourceLevel(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsCreateNotificationsAtActionGroupResourceLevel", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ActionGroupsCreateNotificationsAtActionGroupResourceLevelThenPoll performs ActionGroupsCreateNotificationsAtActionGroupResourceLevel then polls until it's completed
func (c ActionGroupsAPIsClient) ActionGroupsCreateNotificationsAtActionGroupResourceLevelThenPoll(ctx context.Context, id ActionGroupId, input NotificationRequestBody) error {
	result, err := c.ActionGroupsCreateNotificationsAtActionGroupResourceLevel(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ActionGroupsCreateNotificationsAtActionGroupResourceLevel: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ActionGroupsCreateNotificationsAtActionGroupResourceLevel: %+v", err)
	}

	return nil
}

// preparerForActionGroupsCreateNotificationsAtActionGroupResourceLevel prepares the ActionGroupsCreateNotificationsAtActionGroupResourceLevel request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsCreateNotificationsAtActionGroupResourceLevel(ctx context.Context, id ActionGroupId, input NotificationRequestBody) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/createNotifications", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForActionGroupsCreateNotificationsAtActionGroupResourceLevel sends the ActionGroupsCreateNotificationsAtActionGroupResourceLevel request. The method will close the
// http.Response Body if it receives an error.
func (c ActionGroupsAPIsClient) senderForActionGroupsCreateNotificationsAtActionGroupResourceLevel(ctx context.Context, req *http.Request) (future ActionGroupsCreateNotificationsAtActionGroupResourceLevelOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
