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

type TriggerRehydrateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// TriggerRehydrate ...
func (c BackupInstancesClient) TriggerRehydrate(ctx context.Context, id BackupInstanceId, input AzureBackupRehydrationRequest) (result TriggerRehydrateOperationResponse, err error) {
	req, err := c.preparerForTriggerRehydrate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "TriggerRehydrate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForTriggerRehydrate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "TriggerRehydrate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// TriggerRehydrateThenPoll performs TriggerRehydrate then polls until it's completed
func (c BackupInstancesClient) TriggerRehydrateThenPoll(ctx context.Context, id BackupInstanceId, input AzureBackupRehydrationRequest) error {
	result, err := c.TriggerRehydrate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing TriggerRehydrate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after TriggerRehydrate: %+v", err)
	}

	return nil
}

// preparerForTriggerRehydrate prepares the TriggerRehydrate request.
func (c BackupInstancesClient) preparerForTriggerRehydrate(ctx context.Context, id BackupInstanceId, input AzureBackupRehydrationRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/rehydrate", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForTriggerRehydrate sends the TriggerRehydrate request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForTriggerRehydrate(ctx context.Context, req *http.Request) (future TriggerRehydrateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
