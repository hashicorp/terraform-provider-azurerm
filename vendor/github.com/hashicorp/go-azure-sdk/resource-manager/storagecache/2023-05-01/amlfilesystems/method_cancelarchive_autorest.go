package amlfilesystems

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CancelArchiveOperationResponse struct {
	HttpResponse *http.Response
}

// CancelArchive ...
func (c AmlFilesystemsClient) CancelArchive(ctx context.Context, id AmlFilesystemId) (result CancelArchiveOperationResponse, err error) {
	req, err := c.preparerForCancelArchive(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "CancelArchive", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "CancelArchive", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCancelArchive(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "CancelArchive", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCancelArchive prepares the CancelArchive request.
func (c AmlFilesystemsClient) preparerForCancelArchive(ctx context.Context, id AmlFilesystemId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/cancelArchive", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCancelArchive handles the response to the CancelArchive request. The method always
// closes the http.Response Body.
func (c AmlFilesystemsClient) responderForCancelArchive(resp *http.Response) (result CancelArchiveOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
