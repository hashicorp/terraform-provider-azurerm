package virtualmachines

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetRdpFileContentsOperationResponse struct {
	HttpResponse *http.Response
	Model        *RdpConnection
}

// GetRdpFileContents ...
func (c VirtualMachinesClient) GetRdpFileContents(ctx context.Context, id VirtualMachineId) (result GetRdpFileContentsOperationResponse, err error) {
	req, err := c.preparerForGetRdpFileContents(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "GetRdpFileContents", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "GetRdpFileContents", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetRdpFileContents(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "GetRdpFileContents", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetRdpFileContents prepares the GetRdpFileContents request.
func (c VirtualMachinesClient) preparerForGetRdpFileContents(ctx context.Context, id VirtualMachineId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getRdpFileContents", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetRdpFileContents handles the response to the GetRdpFileContents request. The method always
// closes the http.Response Body.
func (c VirtualMachinesClient) responderForGetRdpFileContents(resp *http.Response) (result GetRdpFileContentsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
