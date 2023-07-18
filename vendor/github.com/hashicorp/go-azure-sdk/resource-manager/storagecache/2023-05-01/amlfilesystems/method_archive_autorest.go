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

type ArchiveOperationResponse struct {
	HttpResponse *http.Response
}

// Archive ...
func (c AmlFilesystemsClient) Archive(ctx context.Context, id AmlFilesystemId, input AmlFilesystemArchiveInfo) (result ArchiveOperationResponse, err error) {
	req, err := c.preparerForArchive(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "Archive", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "Archive", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForArchive(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "Archive", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForArchive prepares the Archive request.
func (c AmlFilesystemsClient) preparerForArchive(ctx context.Context, id AmlFilesystemId, input AmlFilesystemArchiveInfo) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/archive", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForArchive handles the response to the Archive request. The method always
// closes the http.Response Body.
func (c AmlFilesystemsClient) responderForArchive(resp *http.Response) (result ArchiveOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
