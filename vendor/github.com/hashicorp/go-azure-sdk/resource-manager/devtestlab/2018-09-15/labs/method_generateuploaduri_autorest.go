package labs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenerateUploadUriOperationResponse struct {
	HttpResponse *http.Response
	Model        *GenerateUploadUriResponse
}

// GenerateUploadUri ...
func (c LabsClient) GenerateUploadUri(ctx context.Context, id LabId, input GenerateUploadUriParameter) (result GenerateUploadUriOperationResponse, err error) {
	req, err := c.preparerForGenerateUploadUri(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "GenerateUploadUri", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "GenerateUploadUri", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGenerateUploadUri(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "GenerateUploadUri", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGenerateUploadUri prepares the GenerateUploadUri request.
func (c LabsClient) preparerForGenerateUploadUri(ctx context.Context, id LabId, input GenerateUploadUriParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/generateUploadUri", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGenerateUploadUri handles the response to the GenerateUploadUri request. The method always
// closes the http.Response Body.
func (c LabsClient) responderForGenerateUploadUri(resp *http.Response) (result GenerateUploadUriOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
