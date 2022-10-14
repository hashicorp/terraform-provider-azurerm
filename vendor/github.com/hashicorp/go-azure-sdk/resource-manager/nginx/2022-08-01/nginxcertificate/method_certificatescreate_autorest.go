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

type CertificatesCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CertificatesCreate ...
func (c NginxCertificateClient) CertificatesCreate(ctx context.Context, id CertificateId, input NginxCertificate) (result CertificatesCreateOperationResponse, err error) {
	req, err := c.preparerForCertificatesCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxcertificate.NginxCertificateClient", "CertificatesCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCertificatesCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxcertificate.NginxCertificateClient", "CertificatesCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CertificatesCreateThenPoll performs CertificatesCreate then polls until it's completed
func (c NginxCertificateClient) CertificatesCreateThenPoll(ctx context.Context, id CertificateId, input NginxCertificate) error {
	result, err := c.CertificatesCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CertificatesCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CertificatesCreate: %+v", err)
	}

	return nil
}

// preparerForCertificatesCreate prepares the CertificatesCreate request.
func (c NginxCertificateClient) preparerForCertificatesCreate(ctx context.Context, id CertificateId, input NginxCertificate) (*http.Request, error) {
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

// senderForCertificatesCreate sends the CertificatesCreate request. The method will close the
// http.Response Body if it receives an error.
func (c NginxCertificateClient) senderForCertificatesCreate(ctx context.Context, req *http.Request) (future CertificatesCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
