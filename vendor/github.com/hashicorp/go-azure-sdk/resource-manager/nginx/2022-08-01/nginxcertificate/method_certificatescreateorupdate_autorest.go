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

type CertificatesCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CertificatesCreateOrUpdate ...
func (c NginxCertificateClient) CertificatesCreateOrUpdate(ctx context.Context, id CertificateId, input NginxCertificate) (result CertificatesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForCertificatesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxcertificate.NginxCertificateClient", "CertificatesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCertificatesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxcertificate.NginxCertificateClient", "CertificatesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CertificatesCreateOrUpdateThenPoll performs CertificatesCreateOrUpdate then polls until it's completed
func (c NginxCertificateClient) CertificatesCreateOrUpdateThenPoll(ctx context.Context, id CertificateId, input NginxCertificate) error {
	result, err := c.CertificatesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CertificatesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CertificatesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForCertificatesCreateOrUpdate prepares the CertificatesCreateOrUpdate request.
func (c NginxCertificateClient) preparerForCertificatesCreateOrUpdate(ctx context.Context, id CertificateId, input NginxCertificate) (*http.Request, error) {
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

// senderForCertificatesCreateOrUpdate sends the CertificatesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c NginxCertificateClient) senderForCertificatesCreateOrUpdate(ctx context.Context, req *http.Request) (future CertificatesCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
