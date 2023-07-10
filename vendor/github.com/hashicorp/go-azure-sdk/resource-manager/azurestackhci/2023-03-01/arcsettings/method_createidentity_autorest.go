package arcsettings

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

type CreateIdentityOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CreateIdentity ...
func (c ArcSettingsClient) CreateIdentity(ctx context.Context, id ArcSettingId) (result CreateIdentityOperationResponse, err error) {
	req, err := c.preparerForCreateIdentity(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "CreateIdentity", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCreateIdentity(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "CreateIdentity", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CreateIdentityThenPoll performs CreateIdentity then polls until it's completed
func (c ArcSettingsClient) CreateIdentityThenPoll(ctx context.Context, id ArcSettingId) error {
	result, err := c.CreateIdentity(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CreateIdentity: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CreateIdentity: %+v", err)
	}

	return nil
}

// preparerForCreateIdentity prepares the CreateIdentity request.
func (c ArcSettingsClient) preparerForCreateIdentity(ctx context.Context, id ArcSettingId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/createArcIdentity", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCreateIdentity sends the CreateIdentity request. The method will close the
// http.Response Body if it receives an error.
func (c ArcSettingsClient) senderForCreateIdentity(ctx context.Context, req *http.Request) (future CreateIdentityOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
