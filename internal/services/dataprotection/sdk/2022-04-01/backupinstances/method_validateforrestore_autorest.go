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

type ValidateForRestoreOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ValidateForRestore ...
func (c BackupInstancesClient) ValidateForRestore(ctx context.Context, id BackupInstanceId, input ValidateRestoreRequestObject) (result ValidateForRestoreOperationResponse, err error) {
	req, err := c.preparerForValidateForRestore(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "ValidateForRestore", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForValidateForRestore(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "ValidateForRestore", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ValidateForRestoreThenPoll performs ValidateForRestore then polls until it's completed
func (c BackupInstancesClient) ValidateForRestoreThenPoll(ctx context.Context, id BackupInstanceId, input ValidateRestoreRequestObject) error {
	result, err := c.ValidateForRestore(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ValidateForRestore: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ValidateForRestore: %+v", err)
	}

	return nil
}

// preparerForValidateForRestore prepares the ValidateForRestore request.
func (c BackupInstancesClient) preparerForValidateForRestore(ctx context.Context, id BackupInstanceId, input ValidateRestoreRequestObject) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/validateRestore", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForValidateForRestore sends the ValidateForRestore request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForValidateForRestore(ctx context.Context, req *http.Request) (future ValidateForRestoreOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
