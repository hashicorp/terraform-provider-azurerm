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

type GenerateCertificateOperationResponse struct {
	HttpResponse *http.Response
	Model        *GenerateCertResponse
}

// GenerateCertificate ...
func (c DevicesClient) GenerateCertificate(ctx context.Context, id DataBoxEdgeDeviceId) (result GenerateCertificateOperationResponse, err error) {
	req, err := c.preparerForGenerateCertificate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GenerateCertificate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GenerateCertificate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGenerateCertificate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "devices.DevicesClient", "GenerateCertificate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGenerateCertificate prepares the GenerateCertificate request.
func (c DevicesClient) preparerForGenerateCertificate(ctx context.Context, id DataBoxEdgeDeviceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/generateCertificate", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGenerateCertificate handles the response to the GenerateCertificate request. The method always
// closes the http.Response Body.
func (c DevicesClient) responderForGenerateCertificate(resp *http.Response) (result GenerateCertificateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
