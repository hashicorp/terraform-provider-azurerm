package azurebackupjob

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

type ExportJobsTriggerOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ExportJobsTrigger ...
func (c AzureBackupJobClient) ExportJobsTrigger(ctx context.Context, id BackupVaultId) (result ExportJobsTriggerOperationResponse, err error) {
	req, err := c.preparerForExportJobsTrigger(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "azurebackupjob.AzureBackupJobClient", "ExportJobsTrigger", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExportJobsTrigger(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "azurebackupjob.AzureBackupJobClient", "ExportJobsTrigger", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExportJobsTriggerThenPoll performs ExportJobsTrigger then polls until it's completed
func (c AzureBackupJobClient) ExportJobsTriggerThenPoll(ctx context.Context, id BackupVaultId) error {
	result, err := c.ExportJobsTrigger(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ExportJobsTrigger: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ExportJobsTrigger: %+v", err)
	}

	return nil
}

// preparerForExportJobsTrigger prepares the ExportJobsTrigger request.
func (c AzureBackupJobClient) preparerForExportJobsTrigger(ctx context.Context, id BackupVaultId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/exportBackupJobs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForExportJobsTrigger sends the ExportJobsTrigger request. The method will close the
// http.Response Body if it receives an error.
func (c AzureBackupJobClient) senderForExportJobsTrigger(ctx context.Context, req *http.Request) (future ExportJobsTriggerOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
