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

type AdhocBackupOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// AdhocBackup ...
func (c BackupInstancesClient) AdhocBackup(ctx context.Context, id BackupInstanceId, input TriggerBackupRequest) (result AdhocBackupOperationResponse, err error) {
	req, err := c.preparerForAdhocBackup(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "AdhocBackup", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForAdhocBackup(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "AdhocBackup", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// AdhocBackupThenPoll performs AdhocBackup then polls until it's completed
func (c BackupInstancesClient) AdhocBackupThenPoll(ctx context.Context, id BackupInstanceId, input TriggerBackupRequest) error {
	result, err := c.AdhocBackup(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AdhocBackup: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after AdhocBackup: %+v", err)
	}

	return nil
}

// preparerForAdhocBackup prepares the AdhocBackup request.
func (c BackupInstancesClient) preparerForAdhocBackup(ctx context.Context, id BackupInstanceId, input TriggerBackupRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/backup", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForAdhocBackup sends the AdhocBackup request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForAdhocBackup(ctx context.Context, req *http.Request) (future AdhocBackupOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
