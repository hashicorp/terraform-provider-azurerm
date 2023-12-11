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

type ResetServicePrincipalProfileOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ResetServicePrincipalProfile ...
func (c ManagedClustersClient) ResetServicePrincipalProfile(ctx context.Context, id commonids.KubernetesClusterId, input ManagedClusterServicePrincipalProfile) (result ResetServicePrincipalProfileOperationResponse, err error) {
	req, err := c.preparerForResetServicePrincipalProfile(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ResetServicePrincipalProfile", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForResetServicePrincipalProfile(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ResetServicePrincipalProfile", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ResetServicePrincipalProfileThenPoll performs ResetServicePrincipalProfile then polls until it's completed
func (c ManagedClustersClient) ResetServicePrincipalProfileThenPoll(ctx context.Context, id commonids.KubernetesClusterId, input ManagedClusterServicePrincipalProfile) error {
	result, err := c.ResetServicePrincipalProfile(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ResetServicePrincipalProfile: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ResetServicePrincipalProfile: %+v", err)
	}

	return nil
}

// preparerForResetServicePrincipalProfile prepares the ResetServicePrincipalProfile request.
func (c ManagedClustersClient) preparerForResetServicePrincipalProfile(ctx context.Context, id commonids.KubernetesClusterId, input ManagedClusterServicePrincipalProfile) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/resetServicePrincipalProfile", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForResetServicePrincipalProfile sends the ResetServicePrincipalProfile request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedClustersClient) senderForResetServicePrincipalProfile(ctx context.Context, req *http.Request) (future ResetServicePrincipalProfileOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
