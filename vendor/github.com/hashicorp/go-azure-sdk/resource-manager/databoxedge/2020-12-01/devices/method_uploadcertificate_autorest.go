package devices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UploadCertificateOperationResponse struct {
	HttpResponse *http.Response
	Model        *UploadCertificateResponse
}

// UploadCertificate ...
func (c DevicesClient) UploadCertificate(ctx context.Context, id DataBoxEdgeDeviceId, input UploadCertificateRequest) (result UploadCertificateOperationResponse, err error) {
	req, err := c.preparerForUploadCertificate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "UploadCertificate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "UploadCertificate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUploadCertificate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "UploadCertificate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUploadCertificate prepares the UploadCertificate request.
func (c DevicesClient) preparerForUploadCertificate(ctx context.Context, id DataBoxEdgeDeviceId, input UploadCertificateRequest) (*http.Request, error) {
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

// responderForUploadCertificate handles the response to the UploadCertificate request. The method always
// closes the http.Response Body.
func (c DevicesClient) responderForUploadCertificate(resp *http.Response) (result UploadCertificateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
