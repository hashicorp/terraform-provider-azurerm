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

type VolumesDeleteReplicationOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// VolumesDeleteReplication ...
func (c VolumesReplicationClient) VolumesDeleteReplication(ctx context.Context, id VolumeId) (result VolumesDeleteReplicationOperationResponse, err error) {
	req, err := c.preparerForVolumesDeleteReplication(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesDeleteReplication", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForVolumesDeleteReplication(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesDeleteReplication", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// VolumesDeleteReplicationThenPoll performs VolumesDeleteReplication then polls until it's completed
func (c VolumesReplicationClient) VolumesDeleteReplicationThenPoll(ctx context.Context, id VolumeId) error {
	result, err := c.VolumesDeleteReplication(ctx, id)
	if err != nil {
		return fmt.Errorf("performing VolumesDeleteReplication: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after VolumesDeleteReplication: %+v", err)
	}

	return nil
}

// preparerForVolumesDeleteReplication prepares the VolumesDeleteReplication request.
func (c VolumesReplicationClient) preparerForVolumesDeleteReplication(ctx context.Context, id VolumeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/deleteReplication", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForVolumesDeleteReplication sends the VolumesDeleteReplication request. The method will close the
// http.Response Body if it receives an error.
func (c VolumesReplicationClient) senderForVolumesDeleteReplication(ctx context.Context, req *http.Request) (future VolumesDeleteReplicationOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
