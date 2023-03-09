package devices

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

type CreateOrUpdateSecuritySettingsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CreateOrUpdateSecuritySettings ...
func (c DevicesClient) CreateOrUpdateSecuritySettings(ctx context.Context, id DataBoxEdgeDeviceId, input SecuritySettings) (result CreateOrUpdateSecuritySettingsOperationResponse, err error) {
	req, err := c.preparerForCreateOrUpdateSecuritySettings(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "CreateOrUpdateSecuritySettings", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCreateOrUpdateSecuritySettings(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "CreateOrUpdateSecuritySettings", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CreateOrUpdateSecuritySettingsThenPoll performs CreateOrUpdateSecuritySettings then polls until it's completed
func (c DevicesClient) CreateOrUpdateSecuritySettingsThenPoll(ctx context.Context, id DataBoxEdgeDeviceId, input SecuritySettings) error {
	result, err := c.CreateOrUpdateSecuritySettings(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CreateOrUpdateSecuritySettings: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CreateOrUpdateSecuritySettings: %+v", err)
	}

	return nil
}

// preparerForCreateOrUpdateSecuritySettings prepares the CreateOrUpdateSecuritySettings request.
func (c DevicesClient) preparerForCreateOrUpdateSecuritySettings(ctx context.Context, id DataBoxEdgeDeviceId, input SecuritySettings) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/securitySettings/default/update", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCreateOrUpdateSecuritySettings sends the CreateOrUpdateSecuritySettings request. The method will close the
// http.Response Body if it receives an error.
func (c DevicesClient) senderForCreateOrUpdateSecuritySettings(ctx context.Context, req *http.Request) (future CreateOrUpdateSecuritySettingsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
