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

type ArcSettingsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ArcSettingsDelete ...
func (c ArcSettingsClient) ArcSettingsDelete(ctx context.Context, id ArcSettingId) (result ArcSettingsDeleteOperationResponse, err error) {
	req, err := c.preparerForArcSettingsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForArcSettingsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "arcsettings.ArcSettingsClient", "ArcSettingsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ArcSettingsDeleteThenPoll performs ArcSettingsDelete then polls until it's completed
func (c ArcSettingsClient) ArcSettingsDeleteThenPoll(ctx context.Context, id ArcSettingId) error {
	result, err := c.ArcSettingsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ArcSettingsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ArcSettingsDelete: %+v", err)
	}

	return nil
}

// preparerForArcSettingsDelete prepares the ArcSettingsDelete request.
func (c ArcSettingsClient) preparerForArcSettingsDelete(ctx context.Context, id ArcSettingId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForArcSettingsDelete sends the ArcSettingsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c ArcSettingsClient) senderForArcSettingsDelete(ctx context.Context, req *http.Request) (future ArcSettingsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
