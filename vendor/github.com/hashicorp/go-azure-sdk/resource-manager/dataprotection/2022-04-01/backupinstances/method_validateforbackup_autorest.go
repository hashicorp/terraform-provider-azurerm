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

type ValidateForBackupOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ValidateForBackup ...
func (c BackupInstancesClient) ValidateForBackup(ctx context.Context, id BackupVaultId, input ValidateForBackupRequest) (result ValidateForBackupOperationResponse, err error) {
	req, err := c.preparerForValidateForBackup(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "ValidateForBackup", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForValidateForBackup(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "ValidateForBackup", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ValidateForBackupThenPoll performs ValidateForBackup then polls until it's completed
func (c BackupInstancesClient) ValidateForBackupThenPoll(ctx context.Context, id BackupVaultId, input ValidateForBackupRequest) error {
	result, err := c.ValidateForBackup(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ValidateForBackup: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ValidateForBackup: %+v", err)
	}

	return nil
}

// preparerForValidateForBackup prepares the ValidateForBackup request.
func (c BackupInstancesClient) preparerForValidateForBackup(ctx context.Context, id BackupVaultId, input ValidateForBackupRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/validateForBackup", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForValidateForBackup sends the ValidateForBackup request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForValidateForBackup(ctx context.Context, req *http.Request) (future ValidateForBackupOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
