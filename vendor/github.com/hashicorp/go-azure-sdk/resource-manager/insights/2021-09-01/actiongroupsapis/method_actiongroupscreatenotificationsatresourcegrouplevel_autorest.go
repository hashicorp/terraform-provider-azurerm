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

type ActionGroupsCreateNotificationsAtResourceGroupLevelOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ActionGroupsCreateNotificationsAtResourceGroupLevel ...
func (c ActionGroupsAPIsClient) ActionGroupsCreateNotificationsAtResourceGroupLevel(ctx context.Context, id commonids.ResourceGroupId, input NotificationRequestBody) (result ActionGroupsCreateNotificationsAtResourceGroupLevelOperationResponse, err error) {
	req, err := c.preparerForActionGroupsCreateNotificationsAtResourceGroupLevel(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsCreateNotificationsAtResourceGroupLevel", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForActionGroupsCreateNotificationsAtResourceGroupLevel(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsCreateNotificationsAtResourceGroupLevel", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ActionGroupsCreateNotificationsAtResourceGroupLevelThenPoll performs ActionGroupsCreateNotificationsAtResourceGroupLevel then polls until it's completed
func (c ActionGroupsAPIsClient) ActionGroupsCreateNotificationsAtResourceGroupLevelThenPoll(ctx context.Context, id commonids.ResourceGroupId, input NotificationRequestBody) error {
	result, err := c.ActionGroupsCreateNotificationsAtResourceGroupLevel(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ActionGroupsCreateNotificationsAtResourceGroupLevel: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ActionGroupsCreateNotificationsAtResourceGroupLevel: %+v", err)
	}

	return nil
}

// preparerForActionGroupsCreateNotificationsAtResourceGroupLevel prepares the ActionGroupsCreateNotificationsAtResourceGroupLevel request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsCreateNotificationsAtResourceGroupLevel(ctx context.Context, id commonids.ResourceGroupId, input NotificationRequestBody) (*http.Request, error) {
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

// senderForActionGroupsCreateNotificationsAtResourceGroupLevel sends the ActionGroupsCreateNotificationsAtResourceGroupLevel request. The method will close the
// http.Response Body if it receives an error.
func (c ActionGroupsAPIsClient) senderForActionGroupsCreateNotificationsAtResourceGroupLevel(ctx context.Context, req *http.Request) (future ActionGroupsCreateNotificationsAtResourceGroupLevelOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
