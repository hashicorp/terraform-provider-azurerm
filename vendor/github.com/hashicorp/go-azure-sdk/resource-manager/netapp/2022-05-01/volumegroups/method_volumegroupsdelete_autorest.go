package volumegroups

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

type VolumeGroupsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// VolumeGroupsDelete ...
func (c VolumeGroupsClient) VolumeGroupsDelete(ctx context.Context, id VolumeGroupId) (result VolumeGroupsDeleteOperationResponse, err error) {
	req, err := c.preparerForVolumeGroupsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForVolumeGroupsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// VolumeGroupsDeleteThenPoll performs VolumeGroupsDelete then polls until it's completed
func (c VolumeGroupsClient) VolumeGroupsDeleteThenPoll(ctx context.Context, id VolumeGroupId) error {
	result, err := c.VolumeGroupsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing VolumeGroupsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after VolumeGroupsDelete: %+v", err)
	}

	return nil
}

// preparerForVolumeGroupsDelete prepares the VolumeGroupsDelete request.
func (c VolumeGroupsClient) preparerForVolumeGroupsDelete(ctx context.Context, id VolumeGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForVolumeGroupsDelete sends the VolumeGroupsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c VolumeGroupsClient) senderForVolumeGroupsDelete(ctx context.Context, req *http.Request) (future VolumeGroupsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
