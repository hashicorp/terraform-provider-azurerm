package blobcontainers

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

type ClearLegalHoldOperationResponse struct {
	HttpResponse *http.Response
	Model        *LegalHold
}

// ClearLegalHold ...
func (c BlobContainersClient) ClearLegalHold(ctx context.Context, id commonids.StorageContainerId, input LegalHold) (result ClearLegalHoldOperationResponse, err error) {
	req, err := c.preparerForClearLegalHold(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "ClearLegalHold", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "ClearLegalHold", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForClearLegalHold(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "blobcontainers.BlobContainersClient", "ClearLegalHold", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForClearLegalHold prepares the ClearLegalHold request.
func (c BlobContainersClient) preparerForClearLegalHold(ctx context.Context, id commonids.StorageContainerId, input LegalHold) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/clearLegalHold", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForClearLegalHold handles the response to the ClearLegalHold request. The method always
// closes the http.Response Body.
func (c BlobContainersClient) responderForClearLegalHold(resp *http.Response) (result ClearLegalHoldOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
