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

type SuspendBackupsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// SuspendBackups ...
func (c BackupInstancesClient) SuspendBackups(ctx context.Context, id BackupInstanceId) (result SuspendBackupsOperationResponse, err error) {
	req, err := c.preparerForSuspendBackups(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "SuspendBackups", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForSuspendBackups(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "SuspendBackups", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// SuspendBackupsThenPoll performs SuspendBackups then polls until it's completed
func (c BackupInstancesClient) SuspendBackupsThenPoll(ctx context.Context, id BackupInstanceId) error {
	result, err := c.SuspendBackups(ctx, id)
	if err != nil {
		return fmt.Errorf("performing SuspendBackups: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after SuspendBackups: %+v", err)
	}

	return nil
}

// preparerForSuspendBackups prepares the SuspendBackups request.
func (c BackupInstancesClient) preparerForSuspendBackups(ctx context.Context, id BackupInstanceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/suspendBackups", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForSuspendBackups sends the SuspendBackups request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForSuspendBackups(ctx context.Context, req *http.Request) (future SuspendBackupsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
