package nginxcertificate

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

type CertificatesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CertificatesDelete ...
func (c NginxCertificateClient) CertificatesDelete(ctx context.Context, id CertificateId) (result CertificatesDeleteOperationResponse, err error) {
	req, err := c.preparerForCertificatesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxcertificate.NginxCertificateClient", "CertificatesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCertificatesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxcertificate.NginxCertificateClient", "CertificatesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CertificatesDeleteThenPoll performs CertificatesDelete then polls until it's completed
func (c NginxCertificateClient) CertificatesDeleteThenPoll(ctx context.Context, id CertificateId) error {
	result, err := c.CertificatesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CertificatesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CertificatesDelete: %+v", err)
	}

	return nil
}

// preparerForCertificatesDelete prepares the CertificatesDelete request.
func (c NginxCertificateClient) preparerForCertificatesDelete(ctx context.Context, id CertificateId) (*http.Request, error) {
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

// senderForCertificatesDelete sends the CertificatesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c NginxCertificateClient) senderForCertificatesDelete(ctx context.Context, req *http.Request) (future CertificatesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
