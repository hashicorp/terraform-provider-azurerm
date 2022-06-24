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

type SyncBackupInstanceOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SyncBackupInstance ...
func (c BackupInstancesClient) SyncBackupInstance(ctx context.Context, id BackupInstanceId, input SyncBackupInstanceRequest) (result SyncBackupInstanceOperationResponse, err error) {
	req, err := c.preparerForSyncBackupInstance(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "SyncBackupInstance", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSyncBackupInstance(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "SyncBackupInstance", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SyncBackupInstanceThenPoll performs SyncBackupInstance then polls until it's completed
func (c BackupInstancesClient) SyncBackupInstanceThenPoll(ctx context.Context, id BackupInstanceId, input SyncBackupInstanceRequest) error {
	result, err := c.SyncBackupInstance(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing SyncBackupInstance: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SyncBackupInstance: %+v", err)
	}

	return nil
}

// preparerForSyncBackupInstance prepares the SyncBackupInstance request.
func (c BackupInstancesClient) preparerForSyncBackupInstance(ctx context.Context, id BackupInstanceId, input SyncBackupInstanceRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/sync", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSyncBackupInstance sends the SyncBackupInstance request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForSyncBackupInstance(ctx context.Context, req *http.Request) (future SyncBackupInstanceOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
