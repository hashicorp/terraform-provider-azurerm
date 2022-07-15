package workspaces

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

type DiagnoseOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Diagnose ...
func (c WorkspacesClient) Diagnose(ctx context.Context, id WorkspaceId, input DiagnoseWorkspaceParameters) (result DiagnoseOperationResponse, err error) {
	req, err := c.preparerForDiagnose(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "Diagnose", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDiagnose(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "Diagnose", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DiagnoseThenPoll performs Diagnose then polls until it's completed
func (c WorkspacesClient) DiagnoseThenPoll(ctx context.Context, id WorkspaceId, input DiagnoseWorkspaceParameters) error {
	result, err := c.Diagnose(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Diagnose: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Diagnose: %+v", err)
	}

	return nil
}

// preparerForDiagnose prepares the Diagnose request.
func (c WorkspacesClient) preparerForDiagnose(ctx context.Context, id WorkspaceId, input DiagnoseWorkspaceParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/diagnose", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDiagnose sends the Diagnose request. The method will close the
// http.Response Body if it receives an error.
func (c WorkspacesClient) senderForDiagnose(ctx context.Context, req *http.Request) (future DiagnoseOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}
	future.Poller, err = polling.NewLongRunningPollerFromResponse(ctx, resp, c.Client)
	return
}
