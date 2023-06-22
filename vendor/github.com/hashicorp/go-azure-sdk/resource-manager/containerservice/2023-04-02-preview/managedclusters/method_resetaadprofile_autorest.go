package managedclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResetAADProfileOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ResetAADProfile ...
func (c ManagedClustersClient) ResetAADProfile(ctx context.Context, id commonids.KubernetesClusterId, input ManagedClusterAADProfile) (result ResetAADProfileOperationResponse, err error) {
	req, err := c.preparerForResetAADProfile(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ResetAADProfile", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForResetAADProfile(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ResetAADProfile", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ResetAADProfileThenPoll performs ResetAADProfile then polls until it's completed
func (c ManagedClustersClient) ResetAADProfileThenPoll(ctx context.Context, id commonids.KubernetesClusterId, input ManagedClusterAADProfile) error {
	result, err := c.ResetAADProfile(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ResetAADProfile: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ResetAADProfile: %+v", err)
	}

	return nil
}

// preparerForResetAADProfile prepares the ResetAADProfile request.
func (c ManagedClustersClient) preparerForResetAADProfile(ctx context.Context, id commonids.KubernetesClusterId, input ManagedClusterAADProfile) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/resetAADProfile", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForResetAADProfile sends the ResetAADProfile request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedClustersClient) senderForResetAADProfile(ctx context.Context, req *http.Request) (future ResetAADProfileOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
