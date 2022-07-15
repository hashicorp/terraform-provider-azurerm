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

type VolumeGroupsCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// VolumeGroupsCreate ...
func (c VolumeGroupsClient) VolumeGroupsCreate(ctx context.Context, id VolumeGroupId, input VolumeGroupDetails) (result VolumeGroupsCreateOperationResponse, err error) {
	req, err := c.preparerForVolumeGroupsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForVolumeGroupsCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "volumegroups.VolumeGroupsClient", "VolumeGroupsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// VolumeGroupsCreateThenPoll performs VolumeGroupsCreate then polls until it's completed
func (c VolumeGroupsClient) VolumeGroupsCreateThenPoll(ctx context.Context, id VolumeGroupId, input VolumeGroupDetails) error {
	result, err := c.VolumeGroupsCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing VolumeGroupsCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after VolumeGroupsCreate: %+v", err)
	}

	return nil
}

// preparerForVolumeGroupsCreate prepares the VolumeGroupsCreate request.
func (c VolumeGroupsClient) preparerForVolumeGroupsCreate(ctx context.Context, id VolumeGroupId, input VolumeGroupDetails) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForVolumeGroupsCreate sends the VolumeGroupsCreate request. The method will close the
// http.Response Body if it receives an error.
func (c VolumeGroupsClient) senderForVolumeGroupsCreate(ctx context.Context, req *http.Request) (future VolumeGroupsCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
