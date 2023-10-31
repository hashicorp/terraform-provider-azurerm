package blobservice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SetServicePropertiesOperationResponse struct {
	HttpResponse *http.Response
	Model        *BlobServiceProperties
}

// SetServiceProperties ...
func (c BlobServiceClient) SetServiceProperties(ctx context.Context, id commonids.StorageAccountId, input BlobServiceProperties) (result SetServicePropertiesOperationResponse, err error) {
	req, err := c.preparerForSetServiceProperties(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobservice.BlobServiceClient", "SetServiceProperties", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobservice.BlobServiceClient", "SetServiceProperties", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSetServiceProperties(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobservice.BlobServiceClient", "SetServiceProperties", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSetServiceProperties prepares the SetServiceProperties request.
func (c BlobServiceClient) preparerForSetServiceProperties(ctx context.Context, id commonids.StorageAccountId, input BlobServiceProperties) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/blobServices/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSetServiceProperties handles the response to the SetServiceProperties request. The method always
// closes the http.Response Body.
func (c BlobServiceClient) responderForSetServiceProperties(resp *http.Response) (result SetServicePropertiesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
