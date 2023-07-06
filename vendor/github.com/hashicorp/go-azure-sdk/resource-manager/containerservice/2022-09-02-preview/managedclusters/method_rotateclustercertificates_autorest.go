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

type RotateClusterCertificatesOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// RotateClusterCertificates ...
func (c ManagedClustersClient) RotateClusterCertificates(ctx context.Context, id commonids.KubernetesClusterId) (result RotateClusterCertificatesOperationResponse, err error) {
	req, err := c.preparerForRotateClusterCertificates(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "RotateClusterCertificates", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForRotateClusterCertificates(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "RotateClusterCertificates", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// RotateClusterCertificatesThenPoll performs RotateClusterCertificates then polls until it's completed
func (c ManagedClustersClient) RotateClusterCertificatesThenPoll(ctx context.Context, id commonids.KubernetesClusterId) error {
	result, err := c.RotateClusterCertificates(ctx, id)
	if err != nil {
		return fmt.Errorf("performing RotateClusterCertificates: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after RotateClusterCertificates: %+v", err)
	}

	return nil
}

// preparerForRotateClusterCertificates prepares the RotateClusterCertificates request.
func (c ManagedClustersClient) preparerForRotateClusterCertificates(ctx context.Context, id commonids.KubernetesClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/rotateClusterCertificates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForRotateClusterCertificates sends the RotateClusterCertificates request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedClustersClient) senderForRotateClusterCertificates(ctx context.Context, req *http.Request) (future RotateClusterCertificatesOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
