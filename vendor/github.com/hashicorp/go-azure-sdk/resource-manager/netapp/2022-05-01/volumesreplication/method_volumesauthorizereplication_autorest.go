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

type VolumesAuthorizeReplicationOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// VolumesAuthorizeReplication ...
func (c VolumesReplicationClient) VolumesAuthorizeReplication(ctx context.Context, id VolumeId, input AuthorizeRequest) (result VolumesAuthorizeReplicationOperationResponse, err error) {
	req, err := c.preparerForVolumesAuthorizeReplication(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesAuthorizeReplication", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForVolumesAuthorizeReplication(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumesreplication.VolumesReplicationClient", "VolumesAuthorizeReplication", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// VolumesAuthorizeReplicationThenPoll performs VolumesAuthorizeReplication then polls until it's completed
func (c VolumesReplicationClient) VolumesAuthorizeReplicationThenPoll(ctx context.Context, id VolumeId, input AuthorizeRequest) error {
	result, err := c.VolumesAuthorizeReplication(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing VolumesAuthorizeReplication: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after VolumesAuthorizeReplication: %+v", err)
	}

	return nil
}

// preparerForVolumesAuthorizeReplication prepares the VolumesAuthorizeReplication request.
func (c VolumesReplicationClient) preparerForVolumesAuthorizeReplication(ctx context.Context, id VolumeId, input AuthorizeRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/authorizeReplication", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForVolumesAuthorizeReplication sends the VolumesAuthorizeReplication request. The method will close the
// http.Response Body if it receives an error.
func (c VolumesReplicationClient) senderForVolumesAuthorizeReplication(ctx context.Context, req *http.Request) (future VolumesAuthorizeReplicationOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
