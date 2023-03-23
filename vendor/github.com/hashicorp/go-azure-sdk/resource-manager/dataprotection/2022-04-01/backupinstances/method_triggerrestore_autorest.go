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

type TriggerRestoreOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TriggerRestore ...
func (c BackupInstancesClient) TriggerRestore(ctx context.Context, id BackupInstanceId, input AzureBackupRestoreRequest) (result TriggerRestoreOperationResponse, err error) {
	req, err := c.preparerForTriggerRestore(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "TriggerRestore", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTriggerRestore(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "TriggerRestore", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TriggerRestoreThenPoll performs TriggerRestore then polls until it's completed
func (c BackupInstancesClient) TriggerRestoreThenPoll(ctx context.Context, id BackupInstanceId, input AzureBackupRestoreRequest) error {
	result, err := c.TriggerRestore(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing TriggerRestore: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TriggerRestore: %+v", err)
	}

	return nil
}

// preparerForTriggerRestore prepares the TriggerRestore request.
func (c BackupInstancesClient) preparerForTriggerRestore(ctx context.Context, id BackupInstanceId, input AzureBackupRestoreRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/restore", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForTriggerRestore sends the TriggerRestore request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForTriggerRestore(ctx context.Context, req *http.Request) (future TriggerRestoreOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
