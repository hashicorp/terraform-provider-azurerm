package volumesreplication

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

type VolumesReestablishReplicationOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// VolumesReestablishReplication ...
func (c VolumesReplicationClient) VolumesReestablishReplication(ctx context.Context, id VolumeId, input ReestablishReplicationRequest) (result VolumesReestablishReplicationOperationResponse, err error) {
	req, err := c.preparerForVolumesReestablishReplication(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesReestablishReplication", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForVolumesReestablishReplication(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesReestablishReplication", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// VolumesReestablishReplicationThenPoll performs VolumesReestablishReplication then polls until it's completed
func (c VolumesReplicationClient) VolumesReestablishReplicationThenPoll(ctx context.Context, id VolumeId, input ReestablishReplicationRequest) error {
	result, err := c.VolumesReestablishReplication(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing VolumesReestablishReplication: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after VolumesReestablishReplication: %+v", err)
	}

	return nil
}

// preparerForVolumesReestablishReplication prepares the VolumesReestablishReplication request.
func (c VolumesReplicationClient) preparerForVolumesReestablishReplication(ctx context.Context, id VolumeId, input ReestablishReplicationRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/reestablishReplication", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForVolumesReestablishReplication sends the VolumesReestablishReplication request. The method will close the
// http.Response Body if it receives an error.
func (c VolumesReplicationClient) senderForVolumesReestablishReplication(ctx context.Context, req *http.Request) (future VolumesReestablishReplicationOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
