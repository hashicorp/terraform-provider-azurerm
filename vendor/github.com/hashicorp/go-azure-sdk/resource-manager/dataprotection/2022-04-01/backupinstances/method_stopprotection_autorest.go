package backupinstances

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

type StopProtectionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// StopProtection ...
func (c BackupInstancesClient) StopProtection(ctx context.Context, id BackupInstanceId) (result StopProtectionOperationResponse, err error) {
	req, err := c.preparerForStopProtection(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "StopProtection", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForStopProtection(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "StopProtection", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// StopProtectionThenPoll performs StopProtection then polls until it's completed
func (c BackupInstancesClient) StopProtectionThenPoll(ctx context.Context, id BackupInstanceId) error {
	result, err := c.StopProtection(ctx, id)
	if err != nil {
		return fmt.Errorf("performing StopProtection: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after StopProtection: %+v", err)
	}

	return nil
}

// preparerForStopProtection prepares the StopProtection request.
func (c BackupInstancesClient) preparerForStopProtection(ctx context.Context, id BackupInstanceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/stopProtection", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForStopProtection sends the StopProtection request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForStopProtection(ctx context.Context, req *http.Request) (future StopProtectionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
