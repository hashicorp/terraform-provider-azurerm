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

type CheckAmlFSSubnetsOperationResponse struct {
	HttpResponse *http.Response
}

// CheckAmlFSSubnets ...
func (c AmlFilesystemsClient) CheckAmlFSSubnets(ctx context.Context, id commonids.SubscriptionId, input AmlFilesystemSubnetInfo) (result CheckAmlFSSubnetsOperationResponse, err error) {
	req, err := c.preparerForCheckAmlFSSubnets(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "CheckAmlFSSubnets", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "CheckAmlFSSubnets", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCheckAmlFSSubnets(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "amlfilesystems.AmlFilesystemsClient", "CheckAmlFSSubnets", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCheckAmlFSSubnets prepares the CheckAmlFSSubnets request.
func (c AmlFilesystemsClient) preparerForCheckAmlFSSubnets(ctx context.Context, id commonids.SubscriptionId, input AmlFilesystemSubnetInfo) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.StorageCache/checkAmlFSSubnets", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCheckAmlFSSubnets handles the response to the CheckAmlFSSubnets request. The method always
// closes the http.Response Body.
func (c AmlFilesystemsClient) responderForCheckAmlFSSubnets(resp *http.Response) (result CheckAmlFSSubnetsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
