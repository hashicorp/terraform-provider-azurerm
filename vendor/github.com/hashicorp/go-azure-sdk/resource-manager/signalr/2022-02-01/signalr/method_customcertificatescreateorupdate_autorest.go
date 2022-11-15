package signalr

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

type CustomCertificatesCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CustomCertificatesCreateOrUpdate ...
func (c SignalRClient) CustomCertificatesCreateOrUpdate(ctx context.Context, id CustomCertificateId, input CustomCertificate) (result CustomCertificatesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForCustomCertificatesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCustomCertificatesCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomCertificatesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CustomCertificatesCreateOrUpdateThenPoll performs CustomCertificatesCreateOrUpdate then polls until it's completed
func (c SignalRClient) CustomCertificatesCreateOrUpdateThenPoll(ctx context.Context, id CustomCertificateId, input CustomCertificate) error {
	result, err := c.CustomCertificatesCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CustomCertificatesCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CustomCertificatesCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForCustomCertificatesCreateOrUpdate prepares the CustomCertificatesCreateOrUpdate request.
func (c SignalRClient) preparerForCustomCertificatesCreateOrUpdate(ctx context.Context, id CustomCertificateId, input CustomCertificate) (*http.Request, error) {
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

// senderForCustomCertificatesCreateOrUpdate sends the CustomCertificatesCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c SignalRClient) senderForCustomCertificatesCreateOrUpdate(ctx context.Context, req *http.Request) (future CustomCertificatesCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
