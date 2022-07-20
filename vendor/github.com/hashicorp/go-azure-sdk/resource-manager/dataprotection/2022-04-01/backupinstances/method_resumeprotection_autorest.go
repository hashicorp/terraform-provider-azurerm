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

type ResumeProtectionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ResumeProtection ...
func (c BackupInstancesClient) ResumeProtection(ctx context.Context, id BackupInstanceId) (result ResumeProtectionOperationResponse, err error) {
	req, err := c.preparerForResumeProtection(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "ResumeProtection", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForResumeProtection(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupinstances.BackupInstancesClient", "ResumeProtection", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ResumeProtectionThenPoll performs ResumeProtection then polls until it's completed
func (c BackupInstancesClient) ResumeProtectionThenPoll(ctx context.Context, id BackupInstanceId) error {
	result, err := c.ResumeProtection(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ResumeProtection: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ResumeProtection: %+v", err)
	}

	return nil
}

// preparerForResumeProtection prepares the ResumeProtection request.
func (c BackupInstancesClient) preparerForResumeProtection(ctx context.Context, id BackupInstanceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/resumeProtection", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForResumeProtection sends the ResumeProtection request. The method will close the
// http.Response Body if it receives an error.
func (c BackupInstancesClient) senderForResumeProtection(ctx context.Context, req *http.Request) (future ResumeProtectionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
