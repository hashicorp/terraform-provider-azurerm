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

type ResumeBackupsOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ResumeBackups ...
func (c BackupInstancesClient) ResumeBackups(ctx context.Context, id BackupInstanceId) (result ResumeBackupsOperationResponse, err error) {
	req, err := c.preparerForResumeBackups(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "ResumeBackups", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForResumeBackups(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "ResumeBackups", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ResumeBackupsThenPoll performs ResumeBackups then polls until it's completed
func (c BackupInstancesClient) ResumeBackupsThenPoll(ctx context.Context, id BackupInstanceId) error {
	result, err := c.ResumeBackups(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ResumeBackups: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ResumeBackups: %+v", err)
	}

	return nil
}

// preparerForResumeBackups prepares the ResumeBackups request.
func (c BackupInstancesClient) preparerForResumeBackups(ctx context.Context, id BackupInstanceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/resumeBackups", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForResumeBackups sends the ResumeBackups request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForResumeBackups(ctx context.Context, req *http.Request) (future ResumeBackupsOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
