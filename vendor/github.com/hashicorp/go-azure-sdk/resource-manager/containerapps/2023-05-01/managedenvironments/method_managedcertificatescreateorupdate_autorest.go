package managedenvironments

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

type ManagedCertificatesCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ManagedCertificatesCreateOrUpdate ...
func (c ManagedEnvironmentsClient) ManagedCertificatesCreateOrUpdate(ctx context.Context, id ManagedCertificateId, input ManagedCertificate) (result ManagedCertificatesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForManagedCertificatesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForManagedCertificatesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ManagedCertificatesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ManagedCertificatesCreateOrUpdateThenPoll performs ManagedCertificatesCreateOrUpdate then polls until it's completed
func (c ManagedEnvironmentsClient) ManagedCertificatesCreateOrUpdateThenPoll(ctx context.Context, id ManagedCertificateId, input ManagedCertificate) error {
	result, err := c.ManagedCertificatesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ManagedCertificatesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ManagedCertificatesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForManagedCertificatesCreateOrUpdate prepares the ManagedCertificatesCreateOrUpdate request.
func (c ManagedEnvironmentsClient) preparerForManagedCertificatesCreateOrUpdate(ctx context.Context, id ManagedCertificateId, input ManagedCertificate) (*http.Request, error) {
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

// senderForManagedCertificatesCreateOrUpdate sends the ManagedCertificatesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ManagedEnvironmentsClient) senderForManagedCertificatesCreateOrUpdate(ctx context.Context, req *http.Request) (future ManagedCertificatesCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
