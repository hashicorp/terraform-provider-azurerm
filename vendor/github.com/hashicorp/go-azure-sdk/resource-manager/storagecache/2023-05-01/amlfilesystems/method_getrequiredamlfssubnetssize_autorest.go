package amlfilesystems

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

type GetRequiredAmlFSSubnetsSizeOperationResponse struct {
	HttpResponse *http.Response
	Model        *RequiredAmlFilesystemSubnetsSize
}

// GetRequiredAmlFSSubnetsSize ...
func (c AmlFilesystemsClient) GetRequiredAmlFSSubnetsSize(ctx context.Context, id commonids.SubscriptionId, input RequiredAmlFilesystemSubnetsSizeInfo) (result GetRequiredAmlFSSubnetsSizeOperationResponse, err error) {
	req, err := c.preparerForGetRequiredAmlFSSubnetsSize(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "GetRequiredAmlFSSubnetsSize", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "GetRequiredAmlFSSubnetsSize", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetRequiredAmlFSSubnetsSize(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "GetRequiredAmlFSSubnetsSize", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetRequiredAmlFSSubnetsSize prepares the GetRequiredAmlFSSubnetsSize request.
func (c AmlFilesystemsClient) preparerForGetRequiredAmlFSSubnetsSize(ctx context.Context, id commonids.SubscriptionId, input RequiredAmlFilesystemSubnetsSizeInfo) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.StorageCache/getRequiredAmlFSSubnetsSize", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetRequiredAmlFSSubnetsSize handles the response to the GetRequiredAmlFSSubnetsSize request. The method always
// closes the http.Response Body.
func (c AmlFilesystemsClient) responderForGetRequiredAmlFSSubnetsSize(resp *http.Response) (result GetRequiredAmlFSSubnetsSizeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
