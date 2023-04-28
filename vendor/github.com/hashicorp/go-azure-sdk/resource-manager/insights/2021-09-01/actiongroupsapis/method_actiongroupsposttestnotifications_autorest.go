package actiongroupsapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionGroupsPostTestNotificationsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ActionGroupsPostTestNotifications ...
func (c ActionGroupsAPIsClient) ActionGroupsPostTestNotifications(ctx context.Context, id commonids.SubscriptionId, input NotificationRequestBody) (result ActionGroupsPostTestNotificationsOperationResponse, err error) {
	req, err := c.preparerForActionGroupsPostTestNotifications(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsPostTestNotifications", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForActionGroupsPostTestNotifications(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsPostTestNotifications", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ActionGroupsPostTestNotificationsThenPoll performs ActionGroupsPostTestNotifications then polls until it's completed
func (c ActionGroupsAPIsClient) ActionGroupsPostTestNotificationsThenPoll(ctx context.Context, id commonids.SubscriptionId, input NotificationRequestBody) error {
	result, err := c.ActionGroupsPostTestNotifications(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ActionGroupsPostTestNotifications: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ActionGroupsPostTestNotifications: %+v", err)
	}

	return nil
}

// preparerForActionGroupsPostTestNotifications prepares the ActionGroupsPostTestNotifications request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsPostTestNotifications(ctx context.Context, id commonids.SubscriptionId, input NotificationRequestBody) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/createNotifications", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForActionGroupsPostTestNotifications sends the ActionGroupsPostTestNotifications request. The method will close the
// http.Response Body if it receives an error.
func (c ActionGroupsAPIsClient) senderForActionGroupsPostTestNotifications(ctx context.Context, req *http.Request) (future ActionGroupsPostTestNotificationsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
