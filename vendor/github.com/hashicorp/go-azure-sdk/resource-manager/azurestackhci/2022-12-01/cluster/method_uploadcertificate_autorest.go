package cluster

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

type UploadCertificateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UploadCertificate ...
func (c ClusterClient) UploadCertificate(ctx context.Context, id ClusterId, input UploadCertificateRequest) (result UploadCertificateOperationResponse, err error) {
	req, err := c.preparerForUploadCertificate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cluster.ClusterClient", "UploadCertificate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUploadCertificate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cluster.ClusterClient", "UploadCertificate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UploadCertificateThenPoll performs UploadCertificate then polls until it's completed
func (c ClusterClient) UploadCertificateThenPoll(ctx context.Context, id ClusterId, input UploadCertificateRequest) error {
	result, err := c.UploadCertificate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing UploadCertificate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UploadCertificate: %+v", err)
	}

	return nil
}

// preparerForUploadCertificate prepares the UploadCertificate request.
func (c ClusterClient) preparerForUploadCertificate(ctx context.Context, id ClusterId, input UploadCertificateRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/uploadCertificate", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUploadCertificate sends the UploadCertificate request. The method will close the
// http.Response Body if it receives an error.
func (c ClusterClient) senderForUploadCertificate(ctx context.Context, req *http.Request) (future UploadCertificateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
